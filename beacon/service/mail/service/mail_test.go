// Package service
// Date:   2026/6/26
// Author: Amu
// Description: unit tests for MailService
package service

import (
	"beacon/service/model"
	"beacon/service/schema"
	testutil "beacon/service/testutil"
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestMailService_MailQuery_Success(t *testing.T) {
	r := testutil.NewFakeMailRepo()
	r.MailQueryFn = func(ctx context.Context) (model.Mail, error) {
		return model.Mail{
			Model:    gorm.Model{ID: 1},
			Server:   "smtp.example.com",
			Port:     587,
			Sender:   "alert@example.com",
			Receiver: "admin@example.com",
		}, nil
	}
	svc := NewMailService(r)

	result, err := svc.MailQuery(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 1 {
		t.Errorf("ID = %d, want 1", result.ID)
	}
	if result.Server != "smtp.example.com" {
		t.Errorf("Server = %q, want %q", result.Server, "smtp.example.com")
	}
	if result.Port != 587 {
		t.Errorf("Port = %d, want 587", result.Port)
	}
	if result.Sender != "alert@example.com" {
		t.Errorf("Sender = %q, want %q", result.Sender, "alert@example.com")
	}
	if result.Receiver != "admin@example.com" {
		t.Errorf("Receiver = %q, want %q", result.Receiver, "admin@example.com")
	}
}

func TestMailService_MailQuery_NotFound(t *testing.T) {
	r := testutil.NewFakeMailRepo()
	r.MailQueryFn = func(ctx context.Context) (model.Mail, error) {
		return model.Mail{}, gorm.ErrRecordNotFound
	}
	svc := NewMailService(r)

	result, err := svc.MailQuery(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 0 {
		t.Errorf("ID = %d, want 0 for not-found", result.ID)
	}
	if result.Server != "" {
		t.Errorf("Server = %q, want empty for not-found", result.Server)
	}
}

func TestMailService_MailQuery_Error(t *testing.T) {
	r := testutil.NewFakeMailRepo()
	r.MailQueryFn = func(ctx context.Context) (model.Mail, error) {
		return model.Mail{}, testutil.ErrTest
	}
	svc := NewMailService(r)

	_, err := svc.MailQuery(context.Background())
	if !errors.Is(err, testutil.ErrTest) {
		t.Fatalf("error = %v, want ErrTest", err)
	}
}

func TestMailService_MailCreate(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(ctx context.Context, args schema.MailCreateArgs) error
		wantErr bool
	}{
		{
			name:    "success",
			fn:      func(ctx context.Context, args schema.MailCreateArgs) error { return nil },
			wantErr: false,
		},
		{
			name:    "repo error",
			fn:      func(ctx context.Context, args schema.MailCreateArgs) error { return testutil.ErrTest },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := testutil.NewFakeMailRepo()
			r.MailCreateFn = tt.fn
			svc := NewMailService(r)

			err := svc.MailCreate(context.Background(), schema.MailCreateArgs{
				Server:   "smtp.example.com",
				Port:     587,
				Sender:   "alert@example.com",
				Password: "secret",
				Receiver: "admin@example.com",
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMailService_MailUpdate(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(ctx context.Context, args schema.MailUpdateArgs) error
		wantErr bool
	}{
		{
			name:    "success",
			fn:      func(ctx context.Context, args schema.MailUpdateArgs) error { return nil },
			wantErr: false,
		},
		{
			name:    "repo error",
			fn:      func(ctx context.Context, args schema.MailUpdateArgs) error { return testutil.ErrTest },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := testutil.NewFakeMailRepo()
			r.MailUpdateFn = tt.fn
			svc := NewMailService(r)

			err := svc.MailUpdate(context.Background(), schema.MailUpdateArgs{
				ID:     1,
				Server: "smtp.example.com",
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMailService_MailDelete(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(ctx context.Context, args schema.MailDeleteArgs) error
		wantErr bool
	}{
		{
			name:    "success",
			fn:      func(ctx context.Context, args schema.MailDeleteArgs) error { return nil },
			wantErr: false,
		},
		{
			name:    "repo error",
			fn:      func(ctx context.Context, args schema.MailDeleteArgs) error { return testutil.ErrTest },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := testutil.NewFakeMailRepo()
			r.MailDeleteFn = tt.fn
			svc := NewMailService(r)

			err := svc.MailDelete(context.Background(), schema.MailDeleteArgs{ID: 1})
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMailService_MailTest(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(ctx context.Context, args schema.MailTestArgs) error
		wantErr bool
	}{
		{
			name:    "success",
			fn:      func(ctx context.Context, args schema.MailTestArgs) error { return nil },
			wantErr: false,
		},
		{
			name:    "repo error",
			fn:      func(ctx context.Context, args schema.MailTestArgs) error { return testutil.ErrTest },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := testutil.NewFakeMailRepo()
			r.MailTestFn = tt.fn
			svc := NewMailService(r)

			err := svc.MailTest(context.Background(), schema.MailTestArgs{Receiver: "admin@example.com"})
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
