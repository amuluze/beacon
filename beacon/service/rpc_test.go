package service

import (
	"errors"
	"testing"

	beaconrpc "beacon/pkg/rpc"
)

func TestNewRPCClientHonorsDisabledControl(t *testing.T) {
	globalTunnel = nil
	client, err := NewRPCClient(&Config{Control: Control{Enable: false}}, nil, nil, nil)
	if err != nil {
		t.Fatalf("NewRPCClient() error = %v", err)
	}
	if _, ok := client.(*beaconrpc.DisabledClient); !ok {
		t.Fatalf("NewRPCClient() type = %T, want *rpc.DisabledClient", client)
	}
	if globalTunnel != nil {
		t.Fatal("disabled control must not create a tunnel")
	}
}

func TestNewRPCClientRejectsWeakProductionJoinToken(t *testing.T) {
	_, err := NewRPCClient(&Config{
		App:     App{Env: productionEnv},
		Control: Control{Enable: true, JoinToken: "change-me", TLS: ControlTLS{Enable: true}},
	}, nil, nil, nil)
	if !errors.Is(err, ErrInsecureControlToken) {
		t.Fatalf("NewRPCClient() error = %v, want ErrInsecureControlToken", err)
	}
}

func TestNewRPCClientRequiresTLSInProduction(t *testing.T) {
	_, err := NewRPCClient(&Config{
		App:     App{Env: productionEnv},
		Control: Control{Enable: true, JoinToken: "a-long-and-random-control-token-123", TLS: ControlTLS{Enable: false}},
	}, nil, nil, nil)
	if !errors.Is(err, ErrControlTLSRequired) {
		t.Fatalf("NewRPCClient() error = %v, want ErrControlTLSRequired", err)
	}
}
