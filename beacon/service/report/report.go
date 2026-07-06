// Package report
// Receives monitoring data pushed from Agents via HTTP and persists to DB.
package report

import (
	"errors"
	"fmt"
	"log/slog"

	"beacon/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Service stores monitoring data pushed from Agents.
type Service struct {
	DB    *database.DB
	Token string
}

var ErrMissingAgentID = errors.New("missing agent_id")

func NewService(db *database.DB, token string) *Service {
	return &Service{DB: db, Token: token}
}

// HandleReport is the HTTP POST handler for Agent monitoring data reports.
func (s *Service) HandleReport(c *fiber.Ctx) error {
	token := c.Get("X-Install-Token")
	if token == "" {
		token = c.Query("token")
	}
	if s.Token != "" && token != s.Token {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	var args rpcSchema.MonitorReportArgs
	if err := c.BodyParser(&args); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.Store(args); err != nil {
		if errors.Is(err, ErrMissingAgentID) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		slog.Error("report store failed", "agent", args.AgentID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "store failed"})
	}

	return c.JSON(fiber.Map{"ok": true})
}

// Store persists a batch of monitoring data from an Agent.
func (s *Service) Store(args rpcSchema.MonitorReportArgs) error {
	agentID := args.AgentID
	if agentID == "" {
		return ErrMissingAgentID
	}

	if err := s.DB.RunInTransaction(func(tx *gorm.DB) error {
		// Host - replace latest for this agent.
		if err := tx.Unscoped().Where("agent_id = ?", agentID).Delete(&model.MonitorHost{}).Error; err != nil {
			return fmt.Errorf("delete host report: %w", err)
		}
		if err := tx.Model(&model.MonitorHost{}).Create(&model.MonitorHost{
			AgentID:         agentID,
			Timestamp:       args.Host.Timestamp,
			Uptime:          args.Host.Uptime,
			Hostname:        args.Host.Hostname,
			Os:              args.Host.Os,
			Platform:        args.Host.Platform,
			PlatformVersion: args.Host.PlatformVersion,
			KernelVersion:   args.Host.KernelVersion,
			KernelArch:      args.Host.KernelArch,
		}).Error; err != nil {
			return fmt.Errorf("create host report: %w", err)
		}

		// CPU - append.
		if err := tx.Model(&model.MonitorCPU{}).Create(&model.MonitorCPU{
			AgentID:    agentID,
			Timestamp:  args.CPU.Timestamp,
			CPUPercent: args.CPU.CPUPercent,
		}).Error; err != nil {
			return fmt.Errorf("create cpu report: %w", err)
		}

		// Memory - append.
		if err := tx.Model(&model.MonitorMemory{}).Create(&model.MonitorMemory{
			AgentID:    agentID,
			Timestamp:  args.Memory.Timestamp,
			MemPercent: args.Memory.MemPercent,
			MemTotal:   args.Memory.MemTotal,
			MemUsed:    args.Memory.MemUsed,
		}).Error; err != nil {
			return fmt.Errorf("create memory report: %w", err)
		}

		// Disk - append batch.
		if len(args.Disks) > 0 {
			var disks []model.MonitorDisk
			for _, d := range args.Disks {
				disks = append(disks, model.MonitorDisk{
					AgentID:     agentID,
					Timestamp:   d.Timestamp,
					Device:      d.Device,
					DiskPercent: d.DiskPercent,
					DiskTotal:   d.DiskTotal,
					DiskUsed:    d.DiskUsed,
					DiskRead:    d.DiskRead,
					DiskWrite:   d.DiskWrite,
				})
			}
			if err := tx.Model(&model.MonitorDisk{}).Create(&disks).Error; err != nil {
				return fmt.Errorf("create disk reports: %w", err)
			}
		}

		// Net - append batch.
		if len(args.Nets) > 0 {
			var nets []model.MonitorNet
			for _, n := range args.Nets {
				nets = append(nets, model.MonitorNet{
					AgentID:   agentID,
					Timestamp: n.Timestamp,
					Ethernet:  n.Ethernet,
					NetRecv:   n.NetRecv,
					NetSend:   n.NetSend,
				})
			}
			if err := tx.Model(&model.MonitorNet{}).Create(&nets).Error; err != nil {
				return fmt.Errorf("create net reports: %w", err)
			}
		}

		// Docker - replace latest for this agent.
		if err := tx.Unscoped().Where("agent_id = ?", agentID).Delete(&model.MonitorDocker{}).Error; err != nil {
			return fmt.Errorf("delete docker report: %w", err)
		}
		if err := tx.Model(&model.MonitorDocker{}).Create(&model.MonitorDocker{
			AgentID:       agentID,
			Timestamp:     args.Docker.Timestamp,
			DockerVersion: args.Docker.DockerVersion,
			APIVersion:    args.Docker.APIVersion,
			MinAPIVersion: args.Docker.MinAPIVersion,
			GitCommit:     args.Docker.GitCommit,
			GoVersion:     args.Docker.GoVersion,
			Os:            args.Docker.Os,
			Arch:          args.Docker.Arch,
		}).Error; err != nil {
			return fmt.Errorf("create docker report: %w", err)
		}

		// Container - append batch.
		if len(args.Containers) > 0 {
			var containers []model.MonitorContainer
			for _, c := range args.Containers {
				containers = append(containers, model.MonitorContainer{
					AgentID:     agentID,
					Timestamp:   c.Timestamp,
					ContainerID: c.ContainerID,
					Name:        c.Name,
					Image:       c.Image,
					IP:          c.IP,
					Ports:       c.Ports,
					State:       c.State,
					Uptime:      c.Uptime,
					CPUPercent:  c.CPUPercent,
					MemPercent:  c.MemPercent,
					MemUsage:    c.MemUsage,
					MemLimit:    c.MemLimit,
					Labels:      c.Labels,
				})
			}
			if err := tx.Model(&model.MonitorContainer{}).Create(&containers).Error; err != nil {
				return fmt.Errorf("create container reports: %w", err)
			}
		}

		// Image - replace all for this agent when the batch includes image data.
		if len(args.Images) > 0 {
			if err := tx.Unscoped().Where("agent_id = ?", agentID).Delete(&model.MonitorImage{}).Error; err != nil {
				return fmt.Errorf("delete image reports: %w", err)
			}
			var images []model.MonitorImage
			for _, im := range args.Images {
				images = append(images, model.MonitorImage{
					AgentID:   agentID,
					Timestamp: im.Timestamp,
					ImageID:   im.ImageID,
					Name:      im.Name,
					Tag:       im.Tag,
					Created:   im.Created,
					Size:      im.Size,
					Number:    im.Number,
				})
			}
			if err := tx.Model(&model.MonitorImage{}).Create(&images).Error; err != nil {
				return fmt.Errorf("create image reports: %w", err)
			}
		}

		// Network - replace all for this agent when the batch includes network data.
		if len(args.Networks) > 0 {
			if err := tx.Unscoped().Where("agent_id = ?", agentID).Delete(&model.MonitorNetwork{}).Error; err != nil {
				return fmt.Errorf("delete network reports: %w", err)
			}
			var nets []model.MonitorNetwork
			for _, n := range args.Networks {
				nets = append(nets, model.MonitorNetwork{
					AgentID:   agentID,
					Timestamp: n.Timestamp,
					NetworkID: n.NetworkID,
					Name:      n.Name,
					Driver:    n.Driver,
					Scope:     n.Scope,
					Created:   n.Created,
					Internal:  n.Internal,
					Subnet:    n.Subnet,
					Gateway:   n.Gateway,
					Labels:    n.Labels,
				})
			}
			if err := tx.Model(&model.MonitorNetwork{}).Create(&nets).Error; err != nil {
				return fmt.Errorf("create network reports: %w", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	slog.Info("report: stored monitoring data", "agent", agentID,
		"cpu", args.CPU.CPUPercent,
		"mem", args.Memory.MemPercent,
		"disks", len(args.Disks),
		"nets", len(args.Nets),
		"containers", len(args.Containers),
		"images", len(args.Images),
		"networks", len(args.Networks),
	)

	return nil
}
