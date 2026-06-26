// Package agent
package agent

import (
	"context"
	"log/slog"
	"time"

	"amprobe/pkg/errors"
	"amprobe/pkg/fiberx"
	"amprobe/service/model"
	tunnelpkg "common/rpc/tunnel"
	"common/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var Set = wire.NewSet(NewAgentRepo, NewAgentService, NewAgentAPI)

// ── Repository ──

type Repository struct {
	DB *database.DB
}

func NewAgentRepo(db *database.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) List(ctx context.Context) ([]model.Agent, error) {
	var agents []model.Agent
	err := r.DB.Model(&model.Agent{}).Order("last_seen desc").Find(&agents).Error
	return agents, err
}

func (r *Repository) Upsert(agent *model.Agent) error {
	return r.DB.Where("agent_id = ?", agent.AgentID).Assign(agent).FirstOrCreate(agent).Error
}

func (r *Repository) UpdateStatus(agentID string, status string) error {
	return r.DB.Model(&model.Agent{}).Where("agent_id = ?", agentID).Update("status", status).Error
}

func (r *Repository) UpdateLastSeen(agentID string, t time.Time) error {
	return r.DB.Model(&model.Agent{}).Where("agent_id = ?", agentID).Updates(map[string]interface{}{
		"last_seen": t,
		"status":    "online",
	}).Error
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Agent{}).Count(&count).Error
	return count, err
}

// ── Service ──

type Service struct {
	repo   *Repository
	tunnel *tunnelpkg.ServerTunnel
}

func NewAgentService(repo *Repository, tun *tunnelpkg.ServerTunnel) *Service {
	s := &Service{repo: repo, tunnel: tun}
	tun.SetAgentLifecycle(s)
	return s
}

// OnAgentConnect implements tunnel.AgentLifecycle.
func (s *Service) OnAgentConnect(info tunnelpkg.AgentInfo) {
	agent := &model.Agent{
		AgentID:  info.AgentID,
		Version:  info.Version,
		OS:       info.OS,
		Arch:     info.Arch,
		Status:   "online",
		LastSeen: time.Now(),
	}
	if err := s.repo.Upsert(agent); err != nil {
		slog.Error("agent upsert failed", "agent_id", info.AgentID, "err", err)
	}
}

// OnAgentDisconnect implements tunnel.AgentLifecycle.
func (s *Service) OnAgentDisconnect(agentID string) {
	if err := s.repo.UpdateStatus(agentID, "offline"); err != nil {
		slog.Error("agent status update failed", "agent_id", agentID, "err", err)
	}
}

// OnAgentHeartbeat implements tunnel.AgentLifecycle.
func (s *Service) OnAgentHeartbeat(agentID string) {
	if err := s.repo.UpdateLastSeen(agentID, time.Now()); err != nil {
		slog.Debug("agent heartbeat update failed", "agent_id", agentID, "err", err)
	}
}

func (s *Service) List(ctx context.Context) ([]model.Agent, error) {
	return s.repo.List(ctx)
}

func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// ── API ──

type API struct {
	svc *Service
}

func NewAgentAPI(svc *Service) *API {
	return &API{svc: svc}
}

func (a *API) List(c *fiber.Ctx) error {
	agents, err := a.svc.List(c.UserContext())
	if err != nil {
		return fiberx.Failure(c, errors.New400Error(err.Error()))
	}
	return fiberx.Success(c, agents)
}
