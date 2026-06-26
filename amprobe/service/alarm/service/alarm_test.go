// Package service
// Date:   2026/6/26
// Author: Amu
// Description: unit tests for AlarmService
package service

import (
	"amprobe/service/model"
	"amprobe/service/schema"
	testutil "amprobe/service/testutil"
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestAlarmService_AlarmQuery_Success(t *testing.T) {
	r := testutil.NewFakeAlarmRepo()
	r.AlarmQueryFn = func(ctx context.Context) ([]model.AlarmThreshold, error) {
		return []model.AlarmThreshold{
			{Model: gorm.Model{ID: 1}, Type: "cpu", Duration: 60, Threshold: 90},
			{Model: gorm.Model{ID: 2}, Type: "mem", Duration: 120, Threshold: 85},
		}, nil
	}
	svc := NewAlarmService(r)

	result, err := svc.AlarmQuery(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 2 {
		t.Fatalf("len(Data) = %d, want 2", len(result.Data))
	}

	// verify model→schema conversion
	first := result.Data[0]
	if first.ID != 1 {
		t.Errorf("Data[0].ID = %d, want 1", first.ID)
	}
	if first.Type != "cpu" {
		t.Errorf("Data[0].Type = %q, want %q", first.Type, "cpu")
	}
	if first.Duration != 60 {
		t.Errorf("Data[0].Duration = %d, want 60", first.Duration)
	}
	if first.Threshold != 90 {
		t.Errorf("Data[0].Threshold = %d, want 90", first.Threshold)
	}

	second := result.Data[1]
	if second.ID != 2 {
		t.Errorf("Data[1].ID = %d, want 2", second.ID)
	}
	if second.Type != "mem" {
		t.Errorf("Data[1].Type = %q, want %q", second.Type, "mem")
	}
}

func TestAlarmService_AlarmQuery_Empty(t *testing.T) {
	r := testutil.NewFakeAlarmRepo()
	r.AlarmQueryFn = func(ctx context.Context) ([]model.AlarmThreshold, error) {
		return []model.AlarmThreshold{}, nil
	}
	svc := NewAlarmService(r)

	result, err := svc.AlarmQuery(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 0 {
		t.Errorf("len(Data) = %d, want 0", len(result.Data))
	}
}

func TestAlarmService_AlarmQuery_Error(t *testing.T) {
	r := testutil.NewFakeAlarmRepo()
	r.AlarmQueryFn = func(ctx context.Context) ([]model.AlarmThreshold, error) {
		return nil, testutil.ErrTest
	}
	svc := NewAlarmService(r)

	_, err := svc.AlarmQuery(context.Background())
	if !errors.Is(err, testutil.ErrTest) {
		t.Fatalf("error = %v, want ErrTest", err)
	}
}

func TestAlarmService_AlarmUpdate_Success(t *testing.T) {
	r := testutil.NewFakeAlarmRepo()
	r.AlarmUpdateFn = func(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error {
		return nil
	}
	svc := NewAlarmService(r)

	err := svc.AlarmUpdate(context.Background(), schema.AlarmThresholdUpdateArgs{
		ID:        1,
		Type:      "cpu",
		Duration:  120,
		Threshold: 95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAlarmService_AlarmUpdate_Error(t *testing.T) {
	r := testutil.NewFakeAlarmRepo()
	r.AlarmUpdateFn = func(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error {
		return testutil.ErrTest
	}
	svc := NewAlarmService(r)

	err := svc.AlarmUpdate(context.Background(), schema.AlarmThresholdUpdateArgs{ID: 1})
	if !errors.Is(err, testutil.ErrTest) {
		t.Fatalf("error = %v, want ErrTest", err)
	}
}
