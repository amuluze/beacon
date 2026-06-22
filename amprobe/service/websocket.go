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
	rpcSchema "common/rpc/schema"

	"github.com/gofiber/contrib/websocket"
)

type LoggerHandler struct {
	rpcClient *rpc.Client
}

func NewLoggerHandler(client *rpc.Client) *LoggerHandler {
	return &LoggerHandler{rpcClient: client}
}

func (l *LoggerHandler) Handler(c *websocket.Conn) {
	containerID := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	agentID := c.Headers("X-Agent-ID")
	if agentID == "" {
		agentID = c.Query("agent_id")
	}
	if agentID != "" {
		ctx = contextx.NewAgentID(ctx, agentID)
	}

	var reply rpcSchema.ContainerLogsReply
	if err := l.rpcClient.Call(ctx, "ContainerLogs", rpcSchema.ContainerLogsArgs{ContainerID: containerID}, &reply); err != nil {
		slog.Error("read container logs from agent failed", "container_id", containerID, "err", err)
		_ = c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	if len(reply.Data) == 0 {
		return
	}
	if err := c.WriteMessage(websocket.TextMessage, reply.Data); err != nil {
		slog.Error("write container logs websocket message failed", "container_id", containerID, "err", err)
	}
}

type TermHandler struct{}

func NewTermHandler() *TermHandler {
	return &TermHandler{}
}

func (th *TermHandler) Handler(conn *websocket.Conn) {
	const msg = "terminal sessions must be executed by collia agent; server-side ssh execution is disabled"
	_ = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	_ = conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(time.Second))
}
