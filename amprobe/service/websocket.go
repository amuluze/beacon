// Package service
// Date: 2024/3/11 10:38
// Author: Amu
// Description:
package service

import (
	"context"
	"log/slog"
	"time"

	"amprobe/pkg/contextx"
	"amprobe/pkg/rpc"
	"amprobe/service/terminal"
	"common/database"
	rpcSchema "common/rpc/schema"

	"github.com/gofiber/contrib/websocket"
)

type LoggerHandler struct {
	rpcClient rpc.Caller
}

func NewLoggerHandler(client rpc.Caller) *LoggerHandler {
	return &LoggerHandler{rpcClient: client}
}

// NewTerminalHandler creates a terminal handler from service configuration.
func NewTerminalHandler(config *Config, rpcClient rpc.Caller, db *database.DB) *terminal.Handler {
	return terminal.NewHandler(rpcClient, db, config.Session.Directory, config.Session.Enabled)
}

func (l *LoggerHandler) Handler(c *websocket.Conn) {
	containerID := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	agentID := c.Headers("X-Agent-ID")
	if agentID == "" {
		agentID = c.Query("agent_id")
	}
	if agentID != "" {
		ctx = contextx.NewAgentID(ctx, agentID)
	}

	chunkChan, err := l.rpcClient.StreamCall(ctx, "ContainerLogs", rpcSchema.ContainerLogsArgs{ContainerID: containerID})
	if err != nil {
		slog.Error("stream container logs from agent failed", "container_id", containerID, "err", err)
		_ = c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	for chunk := range chunkChan {
		if err := c.WriteMessage(websocket.TextMessage, chunk); err != nil {
			slog.Debug("write container logs websocket message failed", "container_id", containerID, "err", err)
			return
		}
	}
}

type TermHandler struct {
	handler *terminal.Handler
}

// NewTermHandler creates a legacy alias that delegates to terminal.Handler.
func NewTermHandler(handler *terminal.Handler) *TermHandler {
	return &TermHandler{handler: handler}
}

func (th *TermHandler) Handler(conn *websocket.Conn) {
	th.handler.Handle(conn)
}
