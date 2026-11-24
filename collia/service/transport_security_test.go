package service

import (
	"errors"
	"testing"
)

func TestValidateTransportSecurity(t *testing.T) {
	secure := &Config{
		Control: Control{Server: "unix:///data/beacon/control/beacon.sock"},
		Task:    Task{Report: Report{URL: "https://beacon.example.com/api/v1/host/report", Token: "secret"}},
	}
	if err := validateTransportSecurity(secure); err != nil {
		t.Fatalf("secure config rejected: %v", err)
	}

	// 单机部署：非 unix:// 的控制目标一律拒绝，不存在跨主机链路。
	tcpControl := *secure
	tcpControl.Control.Server = "beacon.example.com:17000"
	if err := validateTransportSecurity(&tcpControl); !errors.Is(err, ErrUnixControlTargetRequired) {
		t.Fatalf("control error = %v, want ErrUnixControlTargetRequired", err)
	}

	insecureReport := *secure
	insecureReport.Task.Report.URL = "http://beacon.example.com/api/v1/host/report"
	if err := validateTransportSecurity(&insecureReport); !errors.Is(err, ErrHTTPSReportTargetRequired) {
		t.Fatalf("report error = %v, want ErrHTTPSReportTargetRequired", err)
	}

	emptyControl := *secure
	emptyControl.Control.Server = ""
	if err := validateTransportSecurity(&emptyControl); err != nil {
		t.Fatalf("empty control target must pass (tunnel disabled): %v", err)
	}
}
