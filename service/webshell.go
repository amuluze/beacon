// Package service
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package service

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"log/slog"
	"time"
)

type ShellHandler struct {
	stdinPipe   io.WriteCloser
	comboOutput io.ReadWriter
	session     *ssh.Session
	wsConn      *websocket.Conn
}

type ConnInfo struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type wsWrapper struct {
	*websocket.Conn
}

func (wsw *wsWrapper) Write(p []byte) (n int, err error) {
	writer, err := wsw.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	return writer.Write(p)
}

func (wsw *wsWrapper) Read(p []byte) (n int, err error) {
	for {
		msgType, reader, err := wsw.Conn.NextReader()
		if err != nil {
			return 0, err
		}
		if msgType != websocket.TextMessage {
			continue
		}
		return reader.Read(p)
	}
}

func NewShellHandler() *ShellHandler {
	return &ShellHandler{}
}

func (s *ShellHandler) Handler(c *websocket.Conn) {
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			slog.Error("shell handler read ws message error", "err", err)
			_ = c.WriteMessage(mt, []byte("bad message"))
			return
		}
		fmt.Println("=================", string(msg))
		c.WriteMessage(mt, msg)
	}

	//var cInfo ConnInfo
	//err = json.Unmarshal(msg, &cInfo)
	//if err != nil {
	//	slog.Error("marshal msg error", "err", err)
	//	return
	//}
	//slog.Info("conn info", "info", cInfo)
	//conf := ssh.ClientConfig{
	//	Timeout:         5 * time.Second,
	//	User:            cInfo.Username,
	//	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	//	Auth:            []ssh.AuthMethod{ssh.Password(cInfo.Password)},
	//}
	//
	//sshDial, err := ssh.Dial("tcp", cInfo.IP+":"+cInfo.Port, &conf)
	//if err != nil {
	//	slog.Error("ssh 连接失败", "err", err)
	//	return
	//}
	//defer sshDial.Close()
	//
	//session, err := sshDial.NewSession()
	//if err != nil {
	//	slog.Error("get ssh session error", "err", err)
	//	return
	//}
	//defer session.Close()
	//pipe, _ := session.StdinPipe()
	//rw := io.ReadWriter(&wsWrapper{c})
	//session.Stdout = rw
	//session.Stderr = rw
	//modes := ssh.TerminalModes{
	//	ssh.ECHO:          1,
	//	ssh.TTY_OP_ISPEED: 14400,
	//	ssh.TTY_OP_OSPEED: 14400,
	//}
	////伪造xterm终端
	//err = session.RequestPty("xterm", 100, 100, modes)
	//if err != nil {
	//	slog.Error("第三步:会话伪造终端失败", "err", err)
	//	return
	//}
	//err = session.Shell()
	//if err != nil {
	//	slog.Error("第四步:会话伪造终端失败", "err", err)
	//	return
	//}
	//
	//var h = &ShellHandler{
	//	stdinPipe:   pipe,
	//	comboOutput: rw,
	//	session:     session,
	//	wsConn:      c,
	//}
	////defer session.Close()
	//quitChan := make(chan bool, 3)
	////4.以上初始化信息基本结束.下面是携程读写websocket和ssh管道的操作.也就是信息通信
	//h.start(quitChan)
	////session 等待
	//go h.Wait(quitChan)
	//<-quitChan
}

func (s *ShellHandler) start(quitChan chan bool) {
	go s.receiveMsg(quitChan)
	go s.sendMsg(quitChan)
}

func (s *ShellHandler) Wait(quitChan chan bool) {
	if err := s.session.Wait(); err != nil {
		setQuit(quitChan)
	}
}

func setQuit(quitChan chan bool) {
	quitChan <- true
}

func (s *ShellHandler) sendMsg(quitChan chan bool) {
	wsConn := s.wsConn
	defer setQuit(quitChan)
	ticker := time.NewTicker(time.Millisecond * time.Duration(60))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if s.comboOutput == nil {
				return
			}
			var bytes []byte
			n, _ := s.comboOutput.Read(bytes)
			if n > 0 {
				wsConn.WriteMessage(websocket.TextMessage, bytes)
				//s.comboOutput.buffer.Reset()
			}
		case <-quitChan:
			return
		}

	}
}

// 读取ws信息写入ssh客户端中.
func (s *ShellHandler) receiveMsg(quitChan chan bool) {
	wsConn := s.wsConn
	defer setQuit(quitChan) //告诉其他携程退出
	for {
		select {
		case <-quitChan:
			return
		default:
			//1.websocket 读取信息
			_, data, err := wsConn.ReadMessage()
			fmt.Println("=========readMessae===", string(data))
			if err != nil {
				fmt.Println("receiveWsMsg=>读取ws信息失败", err)
				return
			}
			//2.读取到的数据写入ssh 管道内.
			s.stdinPipe.Write(data)
		}
	}
}
