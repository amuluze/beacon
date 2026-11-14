package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"collia/service/report"
	rpcSchema "common/rpc/schema"
)

// newTestReportClient builds an HTTP report client for the test server.
// The token is empty because the client refuses plain HTTP when a token is set.
func newTestReportClient(url string) *report.Client {
	return report.NewClient(url, "", nil)
}

// fakeDataPusher 模拟隧道数据流：可控制就绪状态与推送结果。
type fakeDataPusher struct {
	ready   bool
	err     error
	methods []string
	payload [][]byte
}

func (f *fakeDataPusher) DataReady() bool { return f.ready }

func (f *fakeDataPusher) PushData(method string, payload []byte) error {
	f.methods = append(f.methods, method)
	f.payload = append(f.payload, payload)
	return f.err
}

func TestPushReport_UsesTunnelDataStream(t *testing.T) {
	pusher := &fakeDataPusher{ready: true}
	task := &TimedTask{agentID: "node-1", dataPusher: pusher}

	args := rpcSchema.MonitorReportArgs{
		AgentID: "node-1",
		CPU:     rpcSchema.CPUReport{CPUPercent: 42},
	}
	task.pushReport(args)

	if len(pusher.methods) != 1 || pusher.methods[0] != rpcSchema.ReportMethod {
		t.Fatalf("tunnel methods = %v, want one MonitorReport", pusher.methods)
	}
	var decoded rpcSchema.MonitorReportArgs
	if err := json.Unmarshal(pusher.payload[0], &decoded); err != nil {
		t.Fatalf("decode pushed payload: %v", err)
	}
	if decoded.AgentID != "node-1" || decoded.CPU.CPUPercent != 42 {
		t.Fatalf("pushed payload = %+v", decoded)
	}
}

func TestPushReport_FallsBackToHTTPWhenDataStreamUnavailable(t *testing.T) {
	var got int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	pusher := &fakeDataPusher{ready: false}
	task := &TimedTask{
		agentID:      "node-1",
		dataPusher:   pusher,
		reportClient: newTestReportClient(server.URL),
	}
	task.pushReport(rpcSchema.MonitorReportArgs{AgentID: "node-1"})

	if got != 1 {
		t.Fatalf("http posts = %d, want 1 fallback push", got)
	}
	if len(pusher.methods) != 0 {
		t.Fatalf("tunnel used while not ready: %v", pusher.methods)
	}
}

func TestPushReport_FallsBackToHTTPWhenTunnelPushFails(t *testing.T) {
	var got int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	pusher := &fakeDataPusher{ready: true, err: errors.New("stream closed")}
	task := &TimedTask{
		agentID:      "node-1",
		dataPusher:   pusher,
		reportClient: newTestReportClient(server.URL),
	}
	task.pushReport(rpcSchema.MonitorReportArgs{AgentID: "node-1"})

	if got != 1 {
		t.Fatalf("http posts = %d, want 1 fallback push", got)
	}
	if len(pusher.methods) != 1 {
		t.Fatalf("tunnel attempts = %v, want 1", pusher.methods)
	}
}

func TestPushReport_DropsBatchWhenNoChannelAvailable(t *testing.T) {
	pusher := &fakeDataPusher{ready: false}
	task := &TimedTask{agentID: "node-1", dataPusher: pusher}

	// 无 HTTP 回落通道且数据流未就绪：丢弃本批，不排队、不阻塞。
	done := make(chan struct{})
	go func() {
		task.pushReport(rpcSchema.MonitorReportArgs{AgentID: "node-1"})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("pushReport blocked on unavailable channel")
	}
	if len(pusher.methods) != 0 {
		t.Fatalf("tunnel attempts = %v, want none", pusher.methods)
	}
}
