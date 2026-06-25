// Package repository
package repository

import (
	"amprobe/service/model"
	"common/database"
	"context"
	"time"

	"github.com/google/wire"
)

var Set = wire.NewSet(AgentRepoSet)

var AgentRepoSet = wire.NewSet(NewAgentRepo, wire.Bind(new(IAgentRepo), new(*AgentRepo)))

type IAgentRepo interface {
	List(ctx context.Context) ([]model.Agent, error)
	GetByID(ctx context.Context, agentID string) (*model.Agent, error)
	Upsert(ctx context.Context, agent *model.Agent) error
	UpdateStatus(ctx context.Context, agentID string, status string) error
	UpdateLastSeen(ctx context.Context, agentID string, t time.Time) error
	Count(ctx context.Context) (int64, error)
}

type AgentRepo struct {
	DB *database.DB
}

func NewAgentRepo(db *database.DB) *AgentRepo {
	return &AgentRepo{DB: db}
}

func (r *AgentRepo) List(ctx context.Context) ([]model.Agent, error) {
	var agents []model.Agent
	err := r.DB.Model(&model.Agent{}).Order("last_seen desc").Find(&agents).Error
	return agents, err
}

func (r *AgentRepo) GetByID(ctx context.Context, agentID string) (*model.Agent, error) {
	var agent model.Agent
	err := r.DB.Where("agent_id = ?", agentID).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *AgentRepo) Upsert(ctx context.Context, agent *model.Agent) error {
	return r.DB.Where("agent_id = ?", agent.AgentID).Assign(agent).FirstOrCreate(agent).Error
}

func (r *AgentRepo) UpdateStatus(ctx context.Context, agentID string, status string) error {
	return r.DB.Model(&model.Agent{}).Where("agent_id = ?", agentID).Update("status", status).Error
}

func (r *AgentRepo) UpdateLastSeen(ctx context.Context, agentID string, t time.Time) error {
	return r.DB.Model(&model.Agent{}).Where("agent_id = ?", agentID).Update("last_seen", t).Error
}

func (r *AgentRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Agent{}).Count(&count).Error
	return count, err
}
