// Package service
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"beacon/service/report"

	rpcSchema "common/rpc/schema"
	tunnelpkg "common/rpc/tunnel"

	"github.com/amuluze/amutool/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(NewInjector)

type Injector struct {
	App           *fiber.App
	Router        *Router
	Config        *Config
	Prepare       *Prepare
	Task          *TimedTask
	ReportService *report.Service
	Logger        *logger.Logger
}

func NewInjector(app *fiber.App, router *Router, prepare *Prepare, config *Config, task *TimedTask, reportService *report.Service, logx *logger.Logger) (*Injector, error) {
	// Monitoring batches are pushed over the control-connection data stream
	// (same connection as the control channel, separate stream). The HTTP
	// report endpoint stays available as a compatibility fallback.
	if tun := ServerTunnelFromHolder(); tun != nil {
		tun.SetDataHandler(newTunnelReportHandler(reportService))
	}
	return &Injector{
		App:           app,
		Router:        router,
		Config:        config,
		Prepare:       prepare,
		Task:          task,
		ReportService: reportService,
		Logger:        logx,
	}, nil
}

// newTunnelReportHandler builds the data-stream handler that stores monitoring
// batches. It reuses the HTTP report service store path, so the batch atomicity
// contract stays identical on both channels.
func newTunnelReportHandler(svc *report.Service) tunnelpkg.DataHandler {
	return func(_ context.Context, agentID string, method string, payload []byte) error {
		if svc == nil {
			return errors.New("report service unavailable")
		}
		if method != rpcSchema.ReportMethod {
			return fmt.Errorf("unsupported tunnel method %q", method)
		}
		var args rpcSchema.MonitorReportArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return fmt.Errorf("decode report args: %w", err)
		}
		// The stream identity comes from the authenticated registration; the
		// payload must not claim a different agent.
		if args.AgentID != agentID {
			return fmt.Errorf("agent id mismatch: stream=%q payload=%q", agentID, args.AgentID)
		}
		return svc.Store(args)
	}
}
