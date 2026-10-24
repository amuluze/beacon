// Package service
// Date:   2026/6/26
// Author: Amu
// Description: unit tests for AlarmService
package service

import (
	"beacon/service/model"
	"beacon/service/schema"
	testutil "beacon/service/testutil"
	"context"
	"errors"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func TestAlarmService_AlarmQuery_Success(t *testing.T) {
	r := testutil.NewFakeAlarmRepo()
	r.AlarmQueryFn = func(ctx context.Context) ([]model.AlarmThreshold, error) {
		return []model.AlarmThreshold{
			{Model: gorm.Model{ID: 1}, Type: "cpu", Duration: 60, Threshold: 90},
			{Model: gorm.Model{ID: 2}, Type: "memory", Duration: 120, Threshold: 85},
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
	if second.Type != "memory" {
		t.Errorf("Data[1].Type = %q, want %q", second.Type, "memory")
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

func TestAlarmService_AlarmMessageSettingQuery(t *testing.T) {
	tests := []struct {
		name  string
		value string
		err   error
		want  string
	}{
		{name: "custom", value: "生产环境主机告警", want: "生产环境主机告警"},
		{name: "missing", err: gorm.ErrRecordNotFound, want: model.DefaultAlarmMessageTitle},
		{name: "invalid stored value", value: "invalid\nheader", want: model.DefaultAlarmMessageTitle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := testutil.NewFakeAlarmRepo()
			r.AlarmMessageSettingQueryFn = func(context.Context) (model.AlarmMessageSetting, error) {
				return model.AlarmMessageSetting{MessageTitle: tt.value}, tt.err
			}
			got, err := NewAlarmService(r).AlarmMessageSettingQuery(context.Background())
			if err != nil {
				t.Fatalf("AlarmMessageSettingQuery() error = %v", err)
			}
			if got.MessageTitle != tt.want {
				t.Fatalf("message title = %q, want %q", got.MessageTitle, tt.want)
			}
		})
	}
}

func TestAlarmService_AlarmMessageSettingQueryPropagatesError(t *testing.T) {
	r := testutil.NewFakeAlarmRepo()
	r.AlarmMessageSettingQueryFn = func(context.Context) (model.AlarmMessageSetting, error) {
		return model.AlarmMessageSetting{}, testutil.ErrTest
	}
	_, err := NewAlarmService(r).AlarmMessageSettingQuery(context.Background())
	if !errors.Is(err, testutil.ErrTest) {
		t.Fatalf("error = %v, want ErrTest", err)
	}
}

func TestAlarmService_AlarmMessageSettingUpdateNormalizesTitle(t *testing.T) {
	r := testutil.NewFakeAlarmRepo()
	var got schema.AlarmMessageSettingUpdateArgs
	r.AlarmMessageSettingUpdateFn = func(_ context.Context, args schema.AlarmMessageSettingUpdateArgs) error {
		got = args
		return nil
	}

	err := NewAlarmService(r).AlarmMessageSettingUpdate(context.Background(), schema.AlarmMessageSettingUpdateArgs{
		MessageTitle: "  生产环境主机告警  ",
	})
	if err != nil {
		t.Fatalf("AlarmMessageSettingUpdate() error = %v", err)
	}
	if got.MessageTitle != "生产环境主机告警" {
		t.Fatalf("persisted title = %q, want trimmed title", got.MessageTitle)
	}
}

func TestAlarmService_AlarmMessageSettingUpdateRejectsInvalidTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{name: "blank", title: "   "},
		{name: "line feed", title: "告警\n标题"},
		{name: "carriage return", title: "告警\r标题"},
		{name: "too long", title: strings.Repeat("告", 129)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := testutil.NewFakeAlarmRepo()
			called := false
			r.AlarmMessageSettingUpdateFn = func(context.Context, schema.AlarmMessageSettingUpdateArgs) error {
				called = true
				return nil
			}
			err := NewAlarmService(r).AlarmMessageSettingUpdate(context.Background(), schema.AlarmMessageSettingUpdateArgs{MessageTitle: tt.title})
			if err == nil {
				t.Fatal("AlarmMessageSettingUpdate() error = nil, want validation error")
			}
			if called {
				t.Fatal("repository update called for invalid title")
			}
		})
	}
}
