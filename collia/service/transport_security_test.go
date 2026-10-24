package service

import (
	"errors"
	"testing"
)

func TestValidateTransportSecurity(t *testing.T) {
	secure := &Config{
		Control: Control{Server: "beacon.example.com:17000", TLS: TLS{Enable: true}},
		Task:    Task{Report: Report{URL: "https://beacon.example.com/api/v1/host/report", Token: "secret"}},
	}
	if err := validateTransportSecurity(secure); err != nil {
		t.Fatalf("secure config rejected: %v", err)
	}

	insecureControl := *secure
	insecureControl.Control.TLS.Enable = false
	if err := validateTransportSecurity(&insecureControl); !errors.Is(err, ErrControlTLSRequired) {
		t.Fatalf("control error = %v, want ErrControlTLSRequired", err)
	}

	insecureReport := *secure
	insecureReport.Task.Report.URL = "http://beacon.example.com/api/v1/host/report"
	if err := validateTransportSecurity(&insecureReport); !errors.Is(err, ErrHTTPSReportTargetRequired) {
		t.Fatalf("report error = %v, want ErrHTTPSReportTargetRequired", err)
	}
}
