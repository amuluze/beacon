// Package health provides HTTP health and readiness probes for the Server.
package health

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// DBPinger is satisfied by any database wrapper that can ping the underlying connection.
type DBPinger interface {
	Ping() error
}

// TunnelMonitor is satisfied by any tunnel manager that can report its own health.
type TunnelMonitor interface {
	AgentCount() int
}

// Probe holds the health check dependencies.
type Probe struct {
	db      DBPinger
	tunnel  TunnelMonitor
	started time.Time
}

// NewProbe creates a new health probe.
func NewProbe() *Probe {
	return &Probe{
		started: time.Now(),
	}
}

// SetDB injects a database pinger for readiness checks.
func (p *Probe) SetDB(db DBPinger) {
	p.db = db
}

// SetTunnel injects a tunnel monitor for readiness checks.
func (p *Probe) SetTunnel(tunnel TunnelMonitor) {
	p.tunnel = tunnel
}

// Liveness returns a simple liveness check (HTTP 200 if the process is running).
func (p *Probe) Liveness(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":    "alive",
		"uptime":    time.Since(p.started).Seconds(),
		"timestamp": time.Now().UTC(),
	})
}

// Readiness returns readiness status. If DB or tunnel dependencies are injected
// and unhealthy, it returns 503 Service Unavailable.
func (p *Probe) Readiness(c *fiber.Ctx) error {
	status := "ready"
	code := http.StatusOK

	if p.db != nil {
		if err := p.db.Ping(); err != nil {
			status = "not_ready"
			code = http.StatusServiceUnavailable
		}
	}
	if p.tunnel != nil && code == http.StatusOK {
		// Tunnel is considered healthy if it has at least zero agents (it is listening).
		// In the future this could be changed to require a minimum agent count.
		_ = p.tunnel.AgentCount()
	}

	return c.Status(code).JSON(fiber.Map{
		"status":    status,
		"timestamp": time.Now().UTC(),
	})
}
