// Package service
// Date:   2026/6/26
// Author: Amu
// Description: unit tests for AuditService
package service

import (
	"beacon/service/model"
	"beacon/service/schema"
	testutil "beacon/service/testutil"
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestAuditService_AuditQuery_Success(t *testing.T) {
	now := time.Now()
	r := testutil.NewFakeAuditRepo()
	r.AuditQueryFn = func(ctx context.Context, args schema.AuditQueryArgs) (model.Audits, error) {
		return model.Audits{
			{Model: gorm.Model{ID: 1, CreatedAt: now}, Username: "admin", Operate: "login"},
			{Model: gorm.Model{ID: 2, CreatedAt: now}, Username: "user1", Operate: "logout"},
		}, nil
	}
	r.AuditCountFn = func(ctx context.Context) (int, error) {
		return 42, nil
	}
	svc := NewAuditService(r)

	args := schema.AuditQueryArgs{Page: 2, Size: 10}
	result, err := svc.AuditQuery(context.Background(), args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// verify pagination fields
	if result.Total != 42 {
		t.Errorf("Total = %d, want 42", result.Total)
	}
	if result.Page != 2 {
		t.Errorf("Page = %d, want 2", result.Page)
	}
	if result.Size != 10 {
		t.Errorf("Size = %d, want 10", result.Size)
	}

	// verify model→schema conversion
	if len(result.Data) != 2 {
		t.Fatalf("len(Data) = %d, want 2", len(result.Data))
	}
	first := result.Data[0]
	if first.ID != 1 {
		t.Errorf("Data[0].ID = %d, want 1", first.ID)
	}
	if first.Username != "admin" {
		t.Errorf("Data[0].Username = %q, want %q", first.Username, "admin")
	}
	if first.Operate != "login" {
		t.Errorf("Data[0].Operate = %q, want %q", first.Operate, "login")
	}
	wantCreated := now.Format("2006-01-02 15:04:05")
	if first.Created != wantCreated {
		t.Errorf("Data[0].Created = %q, want %q", first.Created, wantCreated)
	}
}

func TestAuditService_AuditQuery_Empty(t *testing.T) {
	r := testutil.NewFakeAuditRepo()
	r.AuditQueryFn = func(ctx context.Context, args schema.AuditQueryArgs) (model.Audits, error) {
		return model.Audits{}, nil
	}
	r.AuditCountFn = func(ctx context.Context) (int, error) {
		return 0, nil
	}
	svc := NewAuditService(r)

	args := schema.AuditQueryArgs{Page: 1, Size: 10}
	result, err := svc.AuditQuery(context.Background(), args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 0 {
		t.Errorf("len(Data) = %d, want 0", len(result.Data))
	}
	if result.Total != 0 {
		t.Errorf("Total = %d, want 0", result.Total)
	}
}

func TestAuditService_AuditQuery_Error(t *testing.T) {
	r := testutil.NewFakeAuditRepo()
	r.AuditQueryFn = func(ctx context.Context, args schema.AuditQueryArgs) (model.Audits, error) {
		return nil, testutil.ErrTest
	}
	r.AuditCountFn = func(ctx context.Context) (int, error) {
		return 0, nil
	}
	svc := NewAuditService(r)

	args := schema.AuditQueryArgs{Page: 1, Size: 10}
	_, err := svc.AuditQuery(context.Background(), args)
	if !errors.Is(err, testutil.ErrTest) {
		t.Fatalf("error = %v, want ErrTest", err)
	}
}
