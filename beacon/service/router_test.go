package service

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"beacon/pkg/authcookie"
	"beacon/service/agent"
	authAPI "beacon/service/auth/api"
	healthapi "beacon/service/health/api"
	hostAPI "beacon/service/host/api"
	"beacon/service/terminal"

	"github.com/gofiber/fiber/v2"
)

func TestTerminalWebSocketRoutePrecedesContainerLogRoute(t *testing.T) {
	app := fiber.New()
	router := &Router{
		config:          &Config{},
		loggerHandler:   NewLoggerHandler(nil),
		terminalHandler: terminal.NewHandler(nil, nil),
	}
	router.registerWebSocketRoutes(app)

	terminalIndex := -1
	containerLogIndex := -1
	for index, route := range app.GetRoutes(true) {
		if route.Method != fiber.MethodGet {
			continue
		}
		switch route.Path {
		case "/ws/terminal":
			if terminalIndex == -1 {
				terminalIndex = index
			}
		case "/ws/:id":
			if containerLogIndex == -1 {
				containerLogIndex = index
			}
		case "/ws":
			t.Fatal("legacy /ws terminal route must not be registered")
		}
	}

	if terminalIndex == -1 || containerLogIndex == -1 {
		t.Fatalf("routes missing: terminal=%d container_logs=%d", terminalIndex, containerLogIndex)
	}
	if terminalIndex >= containerLogIndex {
		t.Fatalf("terminal route index %d must precede dynamic log route index %d", terminalIndex, containerLogIndex)
	}
}

func TestNoAuthRouterIsAnonymousAndReadOnly(t *testing.T) {
	cfg := &Config{}
	cfg.App.NoAuth = true
	cfg.Auth.Enable = true
	cfg.Casbin.Enable = true
	app := fiber.New()
	router := &Router{
		config:          cfg,
		authAPI:         &authAPI.AuthAPI{},
		sessionCookies:  authcookie.NewManager(true, 120, 3600),
		healthProbe:     &healthapi.Probe{},
		loggerHandler:   NewLoggerHandler(nil),
		terminalHandler: terminal.NewHandler(nil, nil),
	}
	router.RegisterAPI(app)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/api/v1/auth/user_info", nil))
	if err != nil {
		t.Fatalf("user_info: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("user_info status = %d, want 200", resp.StatusCode)
	}
	var info struct {
		Username string `json:"username"`
		ReadOnly bool   `json:"read_only"`
		NoAuth   bool   `json:"no_auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("decode user_info: %v", err)
	}
	if info.Username != "readonly" || !info.ReadOnly || !info.NoAuth {
		t.Fatalf("unexpected user_info: %#v", info)
	}

	resp, err = app.Test(httptest.NewRequest(fiber.MethodPost, "/api/v1/user/user_create", nil))
	if err != nil {
		t.Fatalf("user_create: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("user_create status = %d, want 403", resp.StatusCode)
	}
}

func TestAgentUpgradeRouteIsAdminOnly(t *testing.T) {
	app := fiber.New()
	router := &Router{agentAPI: &agent.API{}}
	router.registerAgentListRoute(app)

	found := false
	for _, route := range app.GetRoutes(true) {
		if route.Method == fiber.MethodPost && route.Path == "/agent/upgrade" {
			found = true
			if route.Name != "升级 Collia" {
				t.Fatalf("route name = %q", route.Name)
			}
		}
	}
	if !found {
		t.Fatal("POST /agent/upgrade route not registered")
	}
	if _, ok := adminOnlyResourcesMap["升级 Collia"]; !ok {
		t.Fatal("Collia upgrade route must be admin-only")
	}
}

func TestAgentUpdateSettingRouteIsAdminOnly(t *testing.T) {
	app := fiber.New()
	router := &Router{agentAPI: &agent.API{}}
	router.registerAgentListRoute(app)

	getFound, postFound := false, false
	for _, route := range app.GetRoutes(true) {
		if route.Path != "/agent/update_setting" {
			continue
		}
		switch route.Method {
		case fiber.MethodGet:
			getFound = true
			if route.Name != "查询 Agent 自动更新设置" {
				t.Fatalf("GET route name = %q", route.Name)
			}
		case fiber.MethodPost:
			postFound = true
			if route.Name != "更新 Agent 自动更新设置" {
				t.Fatalf("POST route name = %q", route.Name)
			}
		}
	}
	if !getFound || !postFound {
		t.Fatalf("update_setting routes missing: get=%v post=%v", getFound, postFound)
	}
	if _, ok := adminOnlyResourcesMap["更新 Agent 自动更新设置"]; !ok {
		t.Fatal("Agent auto-update setting route must be admin-only")
	}
}

func TestHostRoutesExcludeRemovedSystemSettingsAPI(t *testing.T) {
	app := fiber.New()
	router := &Router{hostAPI: &hostAPI.HostAPI{}}
	router.registerHostRoutes(app)

	registered := make(map[string]struct{})
	for _, route := range app.GetRoutes(true) {
		registered[route.Path] = struct{}{}
	}
	if _, ok := registered["/host/host_info"]; !ok {
		t.Fatal("host monitoring route must remain registered")
	}
	for _, path := range []string{
		"/host/get_dns_settings",
		"/host/set_dns_settings",
		"/host/get_system_time",
		"/host/set_system_time",
		"/host/get_system_timezone_list",
		"/host/get_system_timezone",
		"/host/set_system_timezone",
		"/host/reboot",
		"/host/shutdown",
	} {
		if _, ok := registered[path]; ok {
			t.Errorf("removed system settings route %q is still registered", path)
		}
	}
}
