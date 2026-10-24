package service

import (
	"testing"

	"beacon/service/agent"
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
