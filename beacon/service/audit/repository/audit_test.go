package repository

import (
	"path/filepath"
	"testing"
	"time"

	"beacon/service/model"
	"beacon/service/schema"
	"common/database"
)

func newAuditTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.Audit)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

// Domain Spec I004 — audit repository must support filtering rows by
// agent_id so consumers can pin a single Agent's alerts.
func TestAuditQueryFiltersByAgentID(t *testing.T) {
	db := newAuditTestDB(t)
	now := time.Now()
	if err := db.Create(&[]model.Audit{
		{Username: "system", AgentID: "agent-a", Operate: "[agent-a] cpu alarm"},
		{Username: "system", AgentID: "agent-b", Operate: "[agent-b] cpu alarm"},
		{Username: "system", AgentID: "agent-a", Operate: "[agent-a] memory alarm"},
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	_ = now

	repo := NewAuditRepo(db)

	// Empty agent_id must NOT filter anything — all three rows return.
	all, err := repo.AuditQuery(t.Context(), schema.AuditQueryArgs{
		Type: "system", Page: 1, Size: 20,
	})
	if err != nil {
		t.Fatalf("query all: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("unfiltered count = %d, want 3", len(all))
	}

	// agent_id = "agent-a" must return exactly 2 rows.
	a, err := repo.AuditQuery(t.Context(), schema.AuditQueryArgs{
		Type: "system", AgentID: "agent-a", Page: 1, Size: 20,
	})
	if err != nil {
		t.Fatalf("query agent-a: %v", err)
	}
	if len(a) != 2 {
		t.Fatalf("agent-a count = %d, want 2", len(a))
	}
	for _, row := range a {
		if row.AgentID != "agent-a" {
			t.Fatalf("row AgentID = %q, want agent-a", row.AgentID)
		}
	}

	// agent_id = "agent-b" must return exactly 1 row.
	b, err := repo.AuditQuery(t.Context(), schema.AuditQueryArgs{
		Type: "system", AgentID: "agent-b", Page: 1, Size: 20,
	})
	if err != nil {
		t.Fatalf("query agent-b: %v", err)
	}
	if len(b) != 1 || b[0].AgentID != "agent-b" {
		t.Fatalf("agent-b rows = %v, want exactly one agent-b", b)
	}

	// agent_id that has no rows returns an empty slice (NOT a 404).
	none, err := repo.AuditQuery(t.Context(), schema.AuditQueryArgs{
		Type: "system", AgentID: "agent-c", Page: 1, Size: 20,
	})
	if err != nil {
		t.Fatalf("query unknown agent: %v", err)
	}
	if len(none) != 0 {
		t.Fatalf("unknown agent count = %d, want 0", len(none))
	}
}

func TestAuditCountUsesSameFiltersAsQuery(t *testing.T) {
	db := newAuditTestDB(t)
	if err := db.Create(&[]model.Audit{
		{Username: "system", AgentID: "agent-a", Operate: "cpu alarm"},
		{Username: "system", AgentID: "agent-a", Operate: "memory alarm"},
		{Username: "system", AgentID: "agent-b", Operate: "disk alarm"},
		{Username: "admin", AgentID: "agent-a", Operate: "login"},
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := NewAuditRepo(db)
	total, err := repo.AuditCount(t.Context(), schema.AuditQueryArgs{
		Type: "system", AgentID: "agent-a", Page: 1, Size: 10,
	})
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if total != 2 {
		t.Fatalf("filtered total = %d, want 2", total)
	}
}
