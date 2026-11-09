// Package agent
package agent

import (
	"context"
	stderrors "errors"
	"log/slog"
	"time"

	"beacon/pkg/contextx"
	"beacon/pkg/errors"
	"beacon/pkg/fiberx"
	"beacon/service/model"
	"beacon/service/schema"
	"common/database"
	tunnelpkg "common/rpc/tunnel"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var Set = wire.NewSet(NewAgentRepo, NewAgentService, NewAgentAPI)

const (
	agentStatusOnline  = "online"
	agentStatusOffline = "offline"
)

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

func (r *Repository) Get(ctx context.Context, agentID string) (model.Agent, error) {
	var agent model.Agent
	err := r.DB.WithContext(ctx).Where("agent_id = ?", agentID).First(&agent).Error
	return agent, err
}

// Upsert persists the registration snapshot. Assign must carry a map: with a
// struct, FirstOrCreate skips the UPDATE for an existing row entirely (gorm
// silently no-ops), which would freeze the agent's version/os/arch forever.
func (r *Repository) Upsert(agent *model.Agent) error {
	return r.DB.Where("agent_id = ?", agent.AgentID).Assign(map[string]any{
		"hostname":  agent.Hostname,
		"os":        agent.OS,
		"arch":      agent.Arch,
		"version":   agent.Version,
		"status":    agent.Status,
		"last_seen": agent.LastSeen,
	}).FirstOrCreate(agent).Error
}

func (r *Repository) UpdateStatus(agentID string, status string) error {
	return r.DB.Model(&model.Agent{}).Where("agent_id = ?", agentID).Update("status", status).Error
}

func (r *Repository) UpdateLastSeen(agentID string, t time.Time) error {
	return r.DB.Model(&model.Agent{}).Where("agent_id = ?", agentID).Updates(map[string]interface{}{
		"last_seen": t,
		"status":    agentStatusOnline,
	}).Error
}

// MarkAllOffline clears connection state left by a previous Server process.
// A persisted online flag cannot prove that the current process owns a live tunnel.
func (r *Repository) MarkAllOffline() error {
	return r.DB.Model(&model.Agent{}).
		Where("status IS NULL OR status <> ?", agentStatusOffline).
		Update("status", agentStatusOffline).Error
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Agent{}).Count(&count).Error
	return count, err
}

// AutoUpdateEnabled reports whether automatic Collia updates are enabled in
// the console. A missing row means the product default: enabled.
func (r *Repository) AutoUpdateEnabled(ctx context.Context) (bool, error) {
	var setting model.AgentSetting
	err := r.DB.WithContext(ctx).Where("key = ?", model.DefaultAgentAutoUpdateSettingKey).First(&setting).Error
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	if err != nil {
		return true, err
	}
	return setting.Value, nil
}

func (r *Repository) SetAutoUpdateEnabled(ctx context.Context, enabled bool) error {
	return r.DB.WithContext(ctx).
		Where("key = ?", model.DefaultAgentAutoUpdateSettingKey).
		Assign(map[string]any{"value": enabled}).
		FirstOrCreate(&model.AgentSetting{Key: model.DefaultAgentAutoUpdateSettingKey, Value: enabled}).Error
}

// ── Service ──

type Service struct {
	repo    *Repository
	tunnel  onlineChecker
	updater *Updater
}

// SetUpdater wires the Server-side Collia update coordinator. It is set before
// the reverse tunnel starts accepting Agent registrations.
func (s *Service) SetUpdater(updater *Updater) {
	s.updater = updater
}

type onlineChecker interface {
	IsOnline(agentID string) bool
}

func NewAgentService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// SetTunnel wires the tunnel lifecycle hooks into the agent service.
func (s *Service) SetTunnel(tun *tunnelpkg.ServerTunnel) {
	if err := s.repo.MarkAllOffline(); err != nil {
		slog.Error("reset stale agent status failed", "err", err)
	}
	if tun == nil {
		s.tunnel = nil
		return
	}
	s.tunnel = tun
	tun.SetAgentLifecycle(s)
}

// OnAgentConnect implements tunnel.AgentLifecycle.
func (s *Service) OnAgentConnect(info tunnelpkg.AgentInfo) {
	agent := &model.Agent{
		AgentID:  info.AgentID,
		Version:  info.Version,
		OS:       info.OS,
		Arch:     info.Arch,
		Status:   agentStatusOnline,
		LastSeen: time.Now(),
	}
	if err := s.repo.Upsert(agent); err != nil {
		slog.Error("agent upsert failed", "agent_id", info.AgentID, "err", err)
		return
	}
	if s.updater != nil {
		s.updater.Schedule(info)
	}
}

// OnAgentDisconnect implements tunnel.AgentLifecycle.
func (s *Service) OnAgentDisconnect(agentID string) {
	if err := s.repo.UpdateStatus(agentID, agentStatusOffline); err != nil {
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
	agents, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range agents {
		agents[i].Status = agentStatusOffline
		if s.tunnel != nil && s.tunnel.IsOnline(agents[i].AgentID) {
			agents[i].Status = agentStatusOnline
		}
	}
	return agents, nil
}

func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// AutoUpdateSetting is the console-visible state of automatic Collia updates.
// DeploymentEnabled mirrors the deployment-level switch (AgentUpdate.Enable);
// the console switch cannot override a deployment that disabled scheduling.
func (s *Service) AutoUpdateSetting(ctx context.Context) (schema.AgentAutoUpdateSetting, error) {
	enabled, err := s.repo.AutoUpdateEnabled(ctx)
	if err != nil {
		return schema.AgentAutoUpdateSetting{}, err
	}
	setting := schema.AgentAutoUpdateSetting{AutoUpdate: enabled}
	if s.updater != nil {
		setting.DeploymentEnabled = s.updater.DeploymentEnabled()
	}
	return setting, nil
}

func (s *Service) UpdateAutoUpdateSetting(ctx context.Context, enabled bool) (schema.AgentAutoUpdateSetting, error) {
	if err := s.repo.SetAutoUpdateEnabled(ctx, enabled); err != nil {
		return schema.AgentAutoUpdateSetting{}, err
	}
	return s.AutoUpdateSetting(ctx)
}

func (s *Service) Upgrade(ctx context.Context, agentID string) (UpdateResult, error) {
	if agentID == "" {
		return UpdateResult{}, contextx.ErrAgentIDRequired
	}
	if !contextx.IsValidAgentID(agentID) {
		return UpdateResult{}, contextx.ErrInvalidAgentID
	}
	if s.updater == nil {
		return UpdateResult{}, ErrUpdaterUnavailable
	}
	agent, err := s.repo.Get(ctx, agentID)
	if err != nil {
		return UpdateResult{}, err
	}
	return s.updater.Trigger(ctx, tunnelpkg.AgentInfo{
		AgentID: agent.AgentID,
		Version: agent.Version,
		OS:      agent.OS,
		Arch:    agent.Arch,
	})
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

func (a *API) UpdateSetting(c *fiber.Ctx) error {
	setting, err := a.svc.AutoUpdateSetting(c.UserContext())
	if err != nil {
		return fiberx.Failure(c, fiberx.ServiceError(err))
	}
	return fiberx.Success(c, setting)
}

func (a *API) SetUpdateSetting(c *fiber.Ctx) error {
	var args schema.AgentAutoUpdateSettingArgs
	if err := fiberx.ParseBody(c, &args); err != nil {
		return fiberx.Failure(c, errors.New400Error(err.Error()))
	}
	if args.AutoUpdate == nil {
		return fiberx.Failure(c, errors.New400Error("auto_update is required"))
	}
	setting, err := a.svc.UpdateAutoUpdateSetting(c.UserContext(), *args.AutoUpdate)
	if err != nil {
		return fiberx.Failure(c, fiberx.ServiceError(err))
	}
	return fiberx.Success(c, setting)
}

type upgradeArgs struct {
	AgentID string `json:"agent_id"`
}

func (a *API) Upgrade(c *fiber.Ctx) error {
	var args upgradeArgs
	if err := fiberx.ParseBody(c, &args); err != nil {
		return fiberx.Failure(c, errors.New400Error(err.Error()))
	}
	result, err := a.svc.Upgrade(c.UserContext(), args.AgentID)
	if err != nil {
		switch {
		case stderrors.Is(err, ErrUpgradeInProgress):
			return fiberx.Failure(c, errors.New409Error(err.Error()))
		case stderrors.Is(err, ErrUnsupportedArch),
			stderrors.Is(err, ErrDowngradeBlocked),
			stderrors.Is(err, ErrInvalidAgentVersion),
			stderrors.Is(err, ErrInvalidTargetVersion):
			return fiberx.Failure(c, errors.New400Error(err.Error()))
		case stderrors.Is(err, ErrUpdaterUnavailable):
			return fiberx.Failure(c, errors.New503Error(err.Error()))
		}
		return fiberx.Failure(c, fiberx.ServiceError(err))
	}
	return fiberx.Success(c, result)
}
