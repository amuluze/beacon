package ws

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/amuluze/docker"
	"github.com/gofiber/contrib/websocket"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	idChannel := make(chan string)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				slog.Error("read websocket message error", "err", err)
				break
			}
			slog.Info("websocket message received", "type", mt, "msg", string(msg))
			idChannel <- string(msg)
		}
	}()

	for {
		select {
		case containerID := <-idChannel:
			reader, err := l.manager.ContainerLogs(ctx, containerID)
			if err != nil {
				return
			}
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				slog.Info("%s container log: %s", containerID, string(scanner.Bytes()[8:]))
				// scanner.Bytes() 前有 8 个字节的 the HEADER part,需要忽略掉
				// https://medium.com/@dhanushgopinath/reading-docker-container-logs-with-golang-docker-engine-api-702233fac044
				if err = c.WriteJSON(map[string]interface{}{"type": "log", "data": string(scanner.Bytes()[8:])}); err != nil {
					continue
				}
			}
		case <-done:
			slog.Info("websocket closed")
			return
		case <-ctx.Done():
			slog.Info("websocket timeout")
			return
		}
	}
}
