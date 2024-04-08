// Package service
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package service

import (
	"bytes"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"log/slog"
	"os/exec"
)

type ShellHandler struct {
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
		if len(msg) == 0 {
			continue
		}
		fmt.Println(string(msg))
		result, err := RunCommand(string(msg))
		if err != nil {
			slog.Error("shell handler run ws command error", "err", err)
			continue
		}
		err = c.WriteMessage(mt, result)
		if err != nil {
			slog.Error("shell handler write ws command error", "err", err)
		}
	}
}

func RunCommand(command string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", command)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Run()
	return buf.Bytes(), err
}

func RunCommandWithBlock(command string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", command)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Start()
	err = cmd.Wait()
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), err
}
