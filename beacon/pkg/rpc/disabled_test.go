package rpc

import (
	"context"
	"errors"
	"testing"
)

func TestDisabledClient(t *testing.T) {
	client := NewDisabledClient()
	if err := client.Call(context.Background(), "HostInfo", nil, nil); !errors.Is(err, ErrControlDisabled) {
		t.Fatalf("Call() error = %v, want ErrControlDisabled", err)
	}
	if _, err := client.StreamCall(context.Background(), "TerminalSession", nil); !errors.Is(err, ErrControlDisabled) {
		t.Fatalf("StreamCall() error = %v, want ErrControlDisabled", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
