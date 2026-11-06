package service

import (
	"encoding/json"
	"os"
	"testing"

	accountAPI "beacon/service/account/api"
	"beacon/service/agent"
	alarmAPI "beacon/service/alarm/api"
	auditAPI "beacon/service/audit/api"
	authAPI "beacon/service/auth/api"
	containerAPI "beacon/service/container/api"
	dingTalkAPI "beacon/service/dingtalk/api"
	healthAPI "beacon/service/health/api"
	hostAPI "beacon/service/host/api"
	mailAPI "beacon/service/mail/api"
	"beacon/service/terminal"
	"github.com/gofiber/fiber/v2"
)

// TestDocumentationRoutes exports the actual registered surface without sending requests.
func TestDocumentationRoutes(t *testing.T) {
	output := os.Getenv("BEACON_DOC_ROUTES")
	if output == "" {
		t.Skip("documentation export not requested")
	}
	app := fiber.New()
	r := &Router{
		config: &Config{}, agentAPI: &agent.API{},
		accountAPI: &accountAPI.AccountAPI{}, alarmAPI: &alarmAPI.AlarmAPI{},
		auditAPI: &auditAPI.AuditAPI{}, authAPI: &authAPI.AuthAPI{},
		containerAPI: &containerAPI.ContainerAPI{}, dingTalkAPI: &dingTalkAPI.DingTalkAPI{},
		hostAPI: &hostAPI.HostAPI{}, mailAPI: &mailAPI.MailAPI{},
		healthProbe: &healthAPI.Probe{}, loggerHandler: NewLoggerHandler(nil),
		terminalHandler: terminal.NewHandler(nil, nil),
	}
	r.RegisterAPI(app)
	routes := [][2]string{}
	for _, route := range app.GetRoutes(true) {
		if route.Method == fiber.MethodHead {
			continue
		} // Fiber adds HEAD for GET.
		routes = append(routes, [2]string{route.Method, route.Path})
	}
	data, err := json.Marshal(routes)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(output, data, 0600); err != nil {
		t.Fatal(err)
	}
}
