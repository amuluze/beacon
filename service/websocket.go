// Package service
// Date: 2024/3/11 10:38
// Author: Amu
// Description:
package service

import (
	"bufio"
	"context"
	"log/slog"
	"time"

	"github.com/amuluze/amprobe/pkg/docker"
	"github.com/docker/docker/api/types/container"
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
	containerId := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reader, err := l.manager.ContainerLogs(ctx, containerId, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: false})
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
