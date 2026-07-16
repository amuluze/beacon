package service

import (
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type writeTestRouter struct{}

func (writeTestRouter) Register(app *fiber.App) error {
	app.Post("/api/v1/statistics/update", func(ctx *fiber.Ctx) error {
		return ctx.SendString(ctx.IP())
	})
	return nil
}

func (writeTestRouter) Prefixes() []string {
	return []string{"/api/"}
}

func TestNewFiberAppUsesForwardedIPOnlyFromLoopbackProxy(t *testing.T) {
	app := NewFiberApp(&Config{App: App{Env: "production"}}, writeTestRouter{})
	config := app.Config()

	if config.ProxyHeader != fiber.HeaderXForwardedFor {
		t.Fatalf("ProxyHeader = %q, want %q", config.ProxyHeader, fiber.HeaderXForwardedFor)
	}
	if !config.EnableTrustedProxyCheck || !config.EnableIPValidation {
		t.Fatalf("trusted proxy/IP validation 未启用: %+v", config)
	}
	if len(config.TrustedProxies) != 2 || config.TrustedProxies[0] != "127.0.0.1" || config.TrustedProxies[1] != "::1" {
		t.Fatalf("TrustedProxies = %#v, want loopback only", config.TrustedProxies)
	}
}

func TestNewFiberAppLimitsWriteRequestsPerForwardedClientIP(t *testing.T) {
	app := NewFiberApp(&Config{App: App{Env: "production"}}, writeTestRouter{})
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- app.Listener(listener)
	}()
	t.Cleanup(func() {
		if err := app.Shutdown(); err != nil {
			t.Errorf("shutdown: %v", err)
		}
		if err := <-serveErr; err != nil {
			t.Errorf("listener: %v", err)
		}
	})

	client := &http.Client{}
	endpoint := "http://" + listener.Addr().String() + "/api/v1/statistics/update"

	request := func(clientIP string) *http.Response {
		t.Helper()
		req, err := http.NewRequest(http.MethodPost, endpoint, nil)
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		// 官网容器内由 Nitro 通过 loopback 代理到 Go，真实客户端地址位于 X-Forwarded-For。
		req.Header.Set(fiber.HeaderXForwardedFor, clientIP)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("app.Test(%s): %v", clientIP, err)
		}
		return resp
	}

	for i := 0; i < writeLimiterMax; i++ {
		resp := request("203.0.113.10")
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			t.Fatalf("read response %d: %v", i, err)
		}
		if resp.StatusCode != fiber.StatusOK || string(body) != "203.0.113.10" {
			t.Fatalf("request %d = status %d body %q", i, resp.StatusCode, body)
		}
	}

	limited := request("203.0.113.10")
	_ = limited.Body.Close()
	if limited.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("request %d = status %d, want %d", writeLimiterMax+1, limited.StatusCode, fiber.StatusTooManyRequests)
	}

	otherClient := request("203.0.113.11")
	if otherClient.StatusCode != fiber.StatusOK {
		_ = otherClient.Body.Close()
		t.Fatalf("other client = status %d, want %d (independent from 203.0.113.10)", otherClient.StatusCode, fiber.StatusOK)
	}
	if err := otherClient.Body.Close(); err != nil {
		t.Fatalf("close other client response: %v", err)
	}
}
