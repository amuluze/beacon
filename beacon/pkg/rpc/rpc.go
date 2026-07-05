// Package rpc
// Date: 2024/6/12 10:30
// Author: Amu
// Description:
package rpc

import (
	"context"
	"sync"

	"beacon/pkg/contextx"
	tunnelpkg "common/rpc/tunnel"
)

const (
	DefaultAgentID = "default"
	DefaultNetwork = "tcp"
)

// Caller is the interface for making RPC calls to agents.
type Caller interface {
	Call(ctx context.Context, method string, args interface{}, reply interface{}) error
	StreamCall(ctx context.Context, method string, args interface{}) (<-chan []byte, error)
	Close() error
}

// TunnelClient wraps ServerTunnel to implement the Caller interface.
type TunnelClient struct {
	tunnel    *tunnelpkg.ServerTunnel
	defaultID string
	mu        sync.Mutex
	started   bool
}

// NewTunnelClient creates a new tunnel-based RPC caller.
func NewTunnelClient(tunnel *tunnelpkg.ServerTunnel, defaultAgentID string) *TunnelClient {
	if defaultAgentID == "" {
		defaultAgentID = DefaultAgentID
	}
	return &TunnelClient{
		tunnel:    tunnel,
		defaultID: defaultAgentID,
	}
}

func (tc *TunnelClient) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	agentID := contextx.FromAgentID(ctx)
	if agentID == "" {
		agentID = tc.defaultID
	}
	return tc.tunnel.Call(ctx, agentID, method, args, reply)
}

// StreamCall sends an RPC call and returns a channel of response chunks.
func (tc *TunnelClient) StreamCall(ctx context.Context, method string, args interface{}) (<-chan []byte, error) {
	agentID := contextx.FromAgentID(ctx)
	if agentID == "" {
		agentID = tc.defaultID
	}
	chunkChan, err := tc.tunnel.StreamCall(ctx, agentID, method, args)
	if err != nil {
		return nil, err
	}
	// Convert tunnel Frame channel to []byte channel
	out := make(chan []byte, 64)
	go func() {
		defer close(out)
		for {
			select {
			case frame, ok := <-chunkChan:
				if !ok {
					return
				}
				out <- frame.Payload
				if frame.Eos || frame.FrameType == tunnelpkg.FrameType_FRAME_STREAM_END {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

func (tc *TunnelClient) Close() error {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	if tc.started {
		tc.tunnel.Stop()
		tc.started = false
	}
	return nil
}
