package service

import (
	"testing"

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
