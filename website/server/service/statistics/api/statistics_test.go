// Package api
// Date: 2026/07/16
// Author: Amu
// Description: tests for StatisticsAPI with mocked service
package api

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"server/service/schema"
)

type mockService struct {
	updateErr error
}

func (m *mockService) StatisticsQuery(context.Context) (schema.StatisticsQueryReply, error) {
	return schema.StatisticsQueryReply{}, nil
}
func (m *mockService) StatisticsUpdate(context.Context, schema.StatisticsUpdateArgs) (schema.StatisticsUpdateReply, error) {
	return schema.StatisticsUpdateReply{}, m.updateErr
}
func (m *mockService) InstallationReport(context.Context, schema.InstallationReportArgs) (schema.InstallationReportReply, error) {
	return schema.InstallationReportReply{}, nil
}

// id=0 应在 validatex 边界被拒绝，返回 400，不会触达 service 层。
func TestStatisticsUpdate_RejectsInvalidID(t *testing.T) {
	api := NewStatisticsAPI(&mockService{})
	app := fiber.New()
	app.Post("/update", api.StatisticsUpdate)

	req := httptest.NewRequest("POST", "/update", strings.NewReader(`{"id":0}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("id=0 应被校验拒绝，状态码 = %d, want 400", resp.StatusCode)
	}
}

func TestStatisticsUpdate_Accepts(t *testing.T) {
	api := NewStatisticsAPI(&mockService{})
	app := fiber.New()
	app.Post("/update", api.StatisticsUpdate)

	req := httptest.NewRequest("POST", "/update", strings.NewReader(`{"id":3}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("合法 id 状态码 = %d, want 200", resp.StatusCode)
	}
}

// InstallationReport 缺少必填 InstallID 应被校验拒绝。
func TestInstallationReport_RejectsMissingInstallID(t *testing.T) {
	api := NewStatisticsAPI(&mockService{})
	app := fiber.New()
	app.Post("/report", api.InstallationReport)

	req := httptest.NewRequest("POST", "/report", strings.NewReader(`{"image":"img"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("缺少 install_id 应被拒绝，状态码 = %d, want 400", resp.StatusCode)
	}
}
