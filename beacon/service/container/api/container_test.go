package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"beacon/service/container/service"
	"beacon/service/schema"

	"github.com/gofiber/fiber/v2"
)

// stubContainerService 只覆写被测 handler 需要的方法；
// 其余方法由嵌入的接口占位，被意外调用时会因 nil 接口而 panic，从而暴露越界调用。
type stubContainerService struct {
	service.IContainerService
	deleteErr  error
	listCalled bool
}

func (s *stubContainerService) ContainerDelete(_ context.Context, _ schema.ContainerDeleteArgs) error {
	return s.deleteErr
}

func (s *stubContainerService) ContainerList(_ context.Context, _ schema.ContainerQueryArgs) (schema.ContainerQueryRely, error) {
	s.listCalled = true
	return schema.ContainerQueryRely{}, nil
}

func postJSON(t *testing.T, app *fiber.App, path, body string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(fiber.MethodPost, path, strings.NewReader(body))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

// 受保护容器的拒绝必须映射为 403，而不是内部故障 500。
func TestContainerRemoveMapsProtectedRejectionToForbidden(t *testing.T) {
	stub := &stubContainerService{
		deleteErr: fmt.Errorf("%w: beacon", service.ErrProtectedContainer),
	}
	app := fiber.New()
	app.Post("/api/v1/container/container_remove", NewContainerAPI(stub).ContainerRemove)

	resp := postJSON(t, app, "/api/v1/container/container_remove", `{"container_id":"0123456789ab"}`)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
	}
	body := readBody(t, resp)
	if !strings.Contains(body, "container is protected") {
		t.Fatalf("body = %q, want protected rejection detail", body)
	}
}

// 非保护失败的删除仍然返回内部故障，不能被拒绝分支吞掉。
func TestContainerRemoveMapsUnexpectedFailureToInternalError(t *testing.T) {
	stub := &stubContainerService{deleteErr: fmt.Errorf("docker daemon unavailable")}
	app := fiber.New()
	app.Post("/api/v1/container/container_remove", NewContainerAPI(stub).ContainerRemove)

	resp := postJSON(t, app, "/api/v1/container/container_remove", `{"container_id":"0123456789ab"}`)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
}

// 删除成功时返回 204，证明成功路径未被错误映射影响。
func TestContainerRemoveReturnsNoContentOnSuccess(t *testing.T) {
	stub := &stubContainerService{}
	app := fiber.New()
	app.Post("/api/v1/container/container_remove", NewContainerAPI(stub).ContainerRemove)

	resp := postJSON(t, app, "/api/v1/container/container_remove", `{"container_id":"0123456789ab"}`)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

// 缺少必填分页参数属于调用方输入错误：必须在 handler 校验阶段返回 400，且不进入 service。
func TestContainerListRejectsMissingPageBeforeServiceCall(t *testing.T) {
	stub := &stubContainerService{}
	app := fiber.New()
	app.Get("/api/v1/container/containers", NewContainerAPI(stub).ContainerList)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/api/v1/container/containers?size=10", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
	if stub.listCalled {
		t.Fatal("container list service must not be called for an invalid query")
	}
}

// 集合外的 state 取值同样在 handler 校验阶段被拒绝。
func TestContainerListRejectsUnknownState(t *testing.T) {
	stub := &stubContainerService{}
	app := fiber.New()
	app.Get("/api/v1/container/containers", NewContainerAPI(stub).ContainerList)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/api/v1/container/containers?page=1&size=10&state=unknown", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
	if stub.listCalled {
		t.Fatal("container list service must not be called for an unknown state")
	}
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	return string(data)
}
