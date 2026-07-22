// Package rpc
// Date: 2024/06/25
// Author: Amu
// Description: RPC dispatcher using a method registry pattern.
// Routes method name + JSON args to registered handlers.
package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	tunnel "common/rpc/tunnel"
)

// HandlerFunc is the generic handler signature for all RPC methods.
type HandlerFunc func(ctx context.Context, payload []byte, streamSender func(*tunnel.Frame)) ([]byte, error)

// Dispatcher routes incoming RPC calls to registered handlers.
type Dispatcher struct {
	handlers map[string]HandlerFunc
}

// NewDispatcher creates a new dispatcher with all standard handlers registered.
func NewDispatcher(svc *Service) *Dispatcher {
	d := &Dispatcher{handlers: make(map[string]HandlerFunc)}
	registerContainerHandlers(d, svc)
	registerFileHandlers(d, svc)
	registerSystemHandlers(d, svc)
	registerTerminalHandlers(d, svc)
	registerLifecycleHandlers(d, svc)
	return d
}

// Call dispatches the method call with JSON-encoded payload.
func (d *Dispatcher) Call(ctx context.Context, method string, payload []byte, streamSender func(*tunnel.Frame)) ([]byte, error) {
	handler, ok := d.handlers[method]
	if !ok {
		return nil, &UnknownMethodError{Method: method}
	}
	return handler(ctx, payload, streamSender)
}

// Register adds a raw handler for the given method name.
func (d *Dispatcher) Register(method string, fn HandlerFunc) {
	d.handlers[method] = fn
}

// RegisterUnary registers a handler that unmarshals args of type A,
// calls fn, and marshals the reply of type R.
func RegisterUnary[A, R any](d *Dispatcher, method string, fn func(ctx context.Context, args A, reply *R) error) {
	d.Register(method, func(ctx context.Context, payload []byte, _ func(*tunnel.Frame)) ([]byte, error) {
		var args A
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, fmt.Errorf("unmarshal args for %s: %w", method, err)
		}
		var reply R
		if err := fn(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	})
}

// RegisterStream registers a handler that unmarshals args of type A
// and calls fn with the streamSender for streaming output.
func RegisterStream[A any](d *Dispatcher, method string, fn func(ctx context.Context, args A, streamSender func(*tunnel.Frame)) error) {
	d.Register(method, func(ctx context.Context, payload []byte, streamSender func(*tunnel.Frame)) ([]byte, error) {
		var args A
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, fmt.Errorf("unmarshal args for %s: %w", method, err)
		}
		if err := fn(ctx, args, streamSender); err != nil {
			return nil, err
		}
		return nil, nil // streaming methods don't return a payload
	})
}

// UnknownMethodError is returned when a method is not recognized.
type UnknownMethodError struct {
	Method string
}

func (e *UnknownMethodError) Error() string {
	return "unknown rpc method: " + e.Method
}
