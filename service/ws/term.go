package ws

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"golang.org/x/crypto/ssh"
)

const (
	MsgData   = '1'
	MsgResize = '2'
)

type Resize struct {
	Columns int
	Rows    int
}

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
		slog.Error("read ssh config error", "err", err)
		return
	}
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
