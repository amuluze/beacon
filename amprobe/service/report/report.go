// Package report
// Receives monitoring data pushed from Agents via HTTP and persists to DB.
package report

import (
	"log/slog"

	"amprobe/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"

	"github.com/gofiber/fiber/v2"
)

// Service stores monitoring data pushed from Agents.
type Service struct {
	DB    *database.DB
	Token string
}

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
		slog.Error("report store failed", "agent", args.AgentID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "store failed"})
	}

	return c.JSON(fiber.Map{"ok": true})
}

// Store persists a batch of monitoring data from an Agent.
func (s *Service) Store(args rpcSchema.MonitorReportArgs) error {
	agentID := args.AgentID

	// Host - replace latest for this agent
	if err := s.DB.Unscoped().Where("agent_id = ?", agentID).Delete(&model.MonitorHost{}).Error; err != nil {
		slog.Error("report: delete host", "agent", agentID, "error", err)
	}
	s.DB.Model(&model.MonitorHost{}).Create(&model.MonitorHost{
		AgentID:         agentID,
		Timestamp:       args.Host.Timestamp,
		Uptime:          args.Host.Uptime,
		Hostname:        args.Host.Hostname,
		Os:              args.Host.Os,
		Platform:        args.Host.Platform,
		PlatformVersion: args.Host.PlatformVersion,
		KernelVersion:   args.Host.KernelVersion,
		KernelArch:      args.Host.KernelArch,
	})

	// CPU - append
	s.DB.Model(&model.MonitorCPU{}).Create(&model.MonitorCPU{
		AgentID:    agentID,
		Timestamp:  args.CPU.Timestamp,
		CPUPercent: args.CPU.CPUPercent,
	})

	// Memory - append
	s.DB.Model(&model.MonitorMemory{}).Create(&model.MonitorMemory{
		AgentID:    agentID,
		Timestamp:  args.Memory.Timestamp,
		MemPercent: args.Memory.MemPercent,
		MemTotal:   args.Memory.MemTotal,
		MemUsed:    args.Memory.MemUsed,
	})

	// Disk - append batch
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
		s.DB.Model(&model.MonitorDisk{}).Create(&disks)
	}

	// Net - append batch
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
		s.DB.Model(&model.MonitorNet{}).Create(&nets)
	}

	// Docker - replace latest for this agent
	if err := s.DB.Unscoped().Where("agent_id = ?", agentID).Delete(&model.MonitorDocker{}).Error; err != nil {
		slog.Error("report: delete docker", "agent", agentID, "error", err)
	}
	s.DB.Model(&model.MonitorDocker{}).Create(&model.MonitorDocker{
		AgentID:       agentID,
		Timestamp:     args.Docker.Timestamp,
		DockerVersion: args.Docker.DockerVersion,
		APIVersion:    args.Docker.APIVersion,
		MinAPIVersion: args.Docker.MinAPIVersion,
		GitCommit:     args.Docker.GitCommit,
		GoVersion:     args.Docker.GoVersion,
		Os:            args.Docker.Os,
		Arch:          args.Docker.Arch,
	})

	// Container - append batch
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
		s.DB.Model(&model.MonitorContainer{}).Create(&containers)
	}

	// Image - replace all for this agent
	if len(args.Images) > 0 {
		if err := s.DB.Unscoped().Where("agent_id = ?", agentID).Delete(&model.MonitorImage{}).Error; err != nil {
			slog.Error("report: delete images", "agent", agentID, "error", err)
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
		s.DB.Model(&model.MonitorImage{}).Create(&images)
	}

	// Network - replace all for this agent
	if len(args.Networks) > 0 {
		if err := s.DB.Unscoped().Where("agent_id = ?", agentID).Delete(&model.MonitorNetwork{}).Error; err != nil {
			slog.Error("report: delete networks", "agent", agentID, "error", err)
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
		s.DB.Model(&model.MonitorNetwork{}).Create(&nets)
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
