package ws

import (
	"context"
	"log/slog"
	"time"

	"github.com/amuluze/docker"
	"github.com/gofiber/contrib/websocket"
)

type ContainerMonitorHandler struct {
	manager *docker.Manager
}

func NewContainerMonitorHandler() *ContainerMonitorHandler {
	manager, err := docker.NewManager()
	if err != nil {
		panic(err)
	}
	return &ContainerMonitorHandler{
		manager: manager,
	}
}

func (l *ContainerMonitorHandler) Handler(c *websocket.Conn) {
	slog.Info("ContainerMonitorHandler started")

	// 创建一个基础上下文和取消函数
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置关闭连接时的清理操作
	defer c.Close()

	// 创建一个定时器，定期获取容器状态
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// 创建一个通道用于接收中断信号
	done := make(chan struct{})
	slog.Info("Client connected to container monitor")

	// 在单独的 goroutine 中处理客户端消息
	go func() {
		defer close(done)
		for {
			mType, msg, err := c.ReadMessage()
			if err != nil {
				slog.Info("Client disconnected", "error", err)
				return
			}
			slog.Info("Client received message", "type", mType, "msg", string(msg))
		}
	}()

	// 主循环，定期发送容器状态
	for {
		select {
		case <-ticker.C:
			// 获取所有容器
			containers, err := l.manager.ListContainer(ctx)
			if err != nil {
				slog.Error("Failed to list containers", "error", err)
				if writeErr := c.WriteJSON(map[string]interface{}{
					"error": "获取容器列表失败: " + err.Error(),
				}); writeErr != nil {
					slog.Error("Failed to write error message", "error", writeErr)
					return
				}
				continue
			}

			// 收集所有容器的状态信息
			stats := make([]map[string]interface{}, 0, len(containers))
			for _, container := range containers {
				// 获取容器的 CPU 和内存使用情况
				cpuPercent, err := l.manager.GetContainerCpu(ctx, container.ID[:6])
				if err != nil {
					slog.Error("Failed to get container CPU stats",
						"container", container.ID[:6],
						"error", err)
					continue
				}

				memPercent, memUsed, memLimit, err := l.manager.GetContainerMem(ctx, container.ID[:6])
				if err != nil {
					slog.Error("Failed to get container memory stats",
						"container", container.ID[:6],
						"error", err)
					continue
				}

				stats = append(stats, map[string]interface{}{
					"id":       container.ID,
					"name":     container.Name,
					"cpu":      cpuPercent,
					"mem":      memPercent,
					"memUsed":  memUsed,
					"memLimit": memLimit,
				})
			}

			// 发送容器状态信息
			if err := c.WriteJSON(map[string]interface{}{
				"type":  "containerStats",
				"stats": stats,
			}); err != nil {
				slog.Error("Failed to write container stats", "error", err)
				return
			}

		case <-done:
			slog.Info("Client connection closed")
			return

		case <-ctx.Done():
			slog.Info("Context cancelled")
			return
		}
	}
}
