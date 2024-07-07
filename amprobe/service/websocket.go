// Package service
// Date: 2024/3/11 10:38
// Author: Amu
// Description:
package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/amuluze/amutool/docker"
	"github.com/amuluze/amutool/errors"
	"github.com/gofiber/contrib/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"log/slog"
	"sync"
	"time"
)

type LoggerHandler struct {
	manager *docker.Manager
}

func NewLoggerHandler() *LoggerHandler {
	manager, err := docker.NewManager()
	if err != nil {
		return nil
	}
	return &LoggerHandler{manager: manager}
}

func (l *LoggerHandler) Handler(c *websocket.Conn) {
	fmt.Println("websocket handler")
	containerId := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reader, err := l.manager.ContainerLogs(ctx, containerId)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(reader)
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			slog.Error("read websocket message error", "err", err)
			_ = c.WriteMessage(mt, []byte("bad message"))
			break
		}
		slog.Info("receive websocket message", "msg", string(msg))
		for scanner.Scan() {
			// scanner.Bytes() 前有 8 个字节的 the HEADER part,需要忽略掉
			// https://medium.com/@dhanushgopinath/reading-docker-container-logs-with-golang-docker-engine-api-702233fac044
			if err = c.WriteMessage(mt, scanner.Bytes()[8:]); err != nil {
				continue
			}
		}
	}
}

const (
	MsgData   = '1'
	MsgResize = '2'
)

type SSHConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type TermHandler struct{}

func NewTermHandler() *TermHandler {
	return &TermHandler{}
}

func (th *TermHandler) Handler(conn *websocket.Conn) {
	var config SSHConfig
	if err := conn.ReadJSON(&config); err != nil {
		return
	}
	slog.Info("ssh config", "config", config)
	conf := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接 ssh 服务器
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), conf)
	if err != nil {
		slog.Error("unable to connect to ssh", "err", err)
		return
	}
	defer client.Close()
	slog.Info("ssh connected")

	term, err := NewTerm(client, conn)
	if err != nil {
		conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(time.Second))
		return
	}
	defer term.Close()

	var logBuffer = bytes.NewBuffer(make([]byte, 0, 10000))
	logBuffer.Reset()

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := term.LoopRead(ctx, logBuffer)
		if err != nil {
			slog.Error("term loop read error: ", "err", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := term.SessionWait()
		if err != nil {
			slog.Error("term session wait error:", "err", err)
		}
		cancel()
	}()
	wg.Wait()
	return
}

type Term struct {
	StdinPipe io.WriteCloser
	Session   *ssh.Session
	Conn      *websocket.Conn
}

func NewTerm(client *ssh.Client, conn *websocket.Conn) (*Term, error) {
	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		slog.Error("unable to create session", "err", err)
		return nil, err
	}
	defer session.Close()
	slog.Info("ssh session created")
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	term := &Term{
		StdinPipe: stdinPipe,
		Session:   session,
		Conn:      conn,
	}
	session.Stdout = term
	session.Stderr = term

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁止回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	}
	if err := session.RequestPty("xterm", 150, 30, modes); err != nil {
		return nil, err
	}
	if err := session.Shell(); err != nil {
		return nil, err
	}
	return term, nil
}

func (t *Term) Close() error {
	if t.Session != nil {
		t.Session.Close()
	}
	return t.Conn.Close()
}

func (t *Term) Read(p []byte) (n int, err error) {
	for {
		msgType, reader, err := t.Conn.NextReader()
		if err != nil {
			return 0, nil
		}
		if msgType != websocket.TextMessage {
			continue
		}
		return reader.Read(p)
	}
}

func (t *Term) LoopRead(ctx context.Context, logBuff *bytes.Buffer) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("LoopRead context done")
		default:
			_, wsData, err := t.Conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("reading websocket message err: %s", err)
			}
			body := t.decode(wsData[1:])
			switch wsData[0] {
			case MsgResize:
				var args Resize
				err := json.Unmarshal(body, &args)
				if err != nil {
					return fmt.Errorf("ssh pty resize windows err: %s", err)
				}
				if args.Columns > 0 && args.Rows > 0 {
					if err := t.Session.WindowChange(args.Rows, args.Columns); err != nil {
						return fmt.Errorf("ssh pty resize windows err: %s", err)
					}
				}
			case MsgData:
				if _, err := t.StdinPipe.Write(body); err != nil {
					return fmt.Errorf("stdinPipe write err: %s", err)
				}
				if _, err := logBuff.Write(body); err != nil {
					return fmt.Errorf("logBuff write err: %s", err)
				}
			}
		}
	}
}

func (t *Term) SessionWait() error {
	if err := t.Session.Wait(); err != nil {
		return err
	}
	return nil
}

func (t *Term) Write(p []byte) (n int, err error) {
	writer, err := t.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()

	return writer.Write(p)
}

func (t *Term) decode(p []byte) []byte {
	decodeString, _ := base64.StdEncoding.DecodeString(string(p))
	return decodeString
}

func (t *Term) encode(p []byte) []byte {
	encodeToString := base64.StdEncoding.EncodeToString(p)
	return []byte(encodeToString)
}

type Resize struct {
	Columns int
	Rows    int
}
