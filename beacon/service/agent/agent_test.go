package agent

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"beacon/service/model"
	"common/database"
	tunnelpkg "common/rpc/tunnel"
)

func newAgentTestDB(t *testing.T) *database.DB {
	t.Helper()

	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "agent")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.Agent)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func createAgent(t *testing.T, db *database.DB, agent model.Agent) {
	t.Helper()
	if err := db.Create(&agent).Error; err != nil {
		t.Fatalf("create agent %q: %v", agent.AgentID, err)
	}
}

func TestServiceListDoesNotExposePersistedOnlineWithoutLiveTunnel(t *testing.T) {
	db := newAgentTestDB(t)
	createAgent(t, db, model.Agent{
		AgentID:  "stale-agent",
		LastSeen: time.Now().Add(-time.Hour),
		Status:   "online",
	})

	service := NewAgentService(NewAgentRepo(db))
	service.SetTunnel(tunnelpkg.NewServerTunnel())
	if err := db.Model(&model.Agent{}).
		Where("agent_id = ?", "stale-agent").
		Update("status", "online").Error; err != nil {
		t.Fatalf("restore stale status fixture: %v", err)
	}

	agents, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list agents: %v", err)
	}
	if len(agents) != 1 {
		t.Fatalf("agent count = %d, want 1", len(agents))
	}
	if got := agents[0].Status; got != "offline" {
		t.Fatalf("status = %q, want offline without a live tunnel", got)
	}
	var stored model.Agent
	if err := db.Where("agent_id = ?", "stale-agent").First(&stored).Error; err != nil {
		t.Fatalf("query stored agent: %v", err)
	}
	if got := stored.Status; got != "online" {
		t.Fatalf("stored status = %q, want online to prove list projection is read-only", got)
	}
}

type stubOnlineChecker map[string]bool

func (s stubOnlineChecker) IsOnline(agentID string) bool {
	return s[agentID]
}

func TestServiceListProjectsLiveTunnelAsOnline(t *testing.T) {
	db := newAgentTestDB(t)
	createAgent(t, db, model.Agent{
		AgentID:  "live-agent",
		LastSeen: time.Now(),
		Status:   "offline",
	})

	service := NewAgentService(NewAgentRepo(db))
	service.tunnel = stubOnlineChecker{"live-agent": true}

	agents, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list agents: %v", err)
	}
	if len(agents) != 1 {
		t.Fatalf("agent count = %d, want 1", len(agents))
	}
	if got := agents[0].Status; got != "online" {
		t.Fatalf("status = %q, want online for a live tunnel", got)
	}
}

func TestSetTunnelResetsPersistedConnectionStatus(t *testing.T) {
	db := newAgentTestDB(t)
	lastSeen := time.Now().Add(-time.Hour)
	createAgent(t, db, model.Agent{
		AgentID:  "stale-agent",
		LastSeen: lastSeen,
		Status:   "online",
	})

	service := NewAgentService(NewAgentRepo(db))
	service.SetTunnel(tunnelpkg.NewServerTunnel())

	var stored model.Agent
	if err := db.Where("agent_id = ?", "stale-agent").First(&stored).Error; err != nil {
		t.Fatalf("query stored agent: %v", err)
	}
	if got := stored.Status; got != "offline" {
		t.Fatalf("stored status = %q, want offline after tunnel initialization", got)
	}
	if !stored.LastSeen.Equal(lastSeen) {
		t.Fatalf("last_seen = %v, want unchanged %v", stored.LastSeen, lastSeen)
	}
}

func TestSetNilTunnelProjectsAgentsOffline(t *testing.T) {
	db := newAgentTestDB(t)
	createAgent(t, db, model.Agent{
		AgentID:  "stale-agent",
		LastSeen: time.Now().Add(-time.Hour),
		Status:   "online",
	})

	service := NewAgentService(NewAgentRepo(db))
	service.SetTunnel(nil)

	agents, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list agents: %v", err)
	}
	if len(agents) != 1 || agents[0].Status != "offline" {
		t.Fatalf("agents = %#v, want one offline agent", agents)
	}
}
