// Package jwtauth
// Date: 2026/6/26
// Author: Amu
// Description: unit tests for JWT auth core logic
package jwtauth

import (
	"sync"
	"testing"
	"time"

	"beacon/pkg/auth"

	"github.com/golang-jwt/jwt/v5"
)

// ── fake Storer (project convention: manual fake, no mock framework) ──

type fakeStore struct {
	mu     sync.Mutex
	data   map[string]bool
	delete bool // if true, Set deletes the key (simulates DestroyToken)
}

func newFakeStore() *fakeStore {
	return &fakeStore{data: make(map[string]bool)}
}

func (f *fakeStore) Set(tokenString string, expiration time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	// DestroyToken sets expiration to 1s; simulate by deleting the key
	if expiration <= time.Second {
		delete(f.data, tokenString)
	} else {
		f.data[tokenString] = true
	}
	return nil
}

func (f *fakeStore) Check(key string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.data[key], nil
}

func (f *fakeStore) Close() error {
	return nil
}

// ── helper ──

func newTestJWTAuth(store Storer) *JWTAuth {
	return New(store, nil,
		SetSigningMethod(jwt.SigningMethodHS256),
		SetSigningKey([]byte("test-key")),
		SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, auth.ErrInvalidToken
			}
			return []byte("test-key"), nil
		}),
		SetExpired(3600),
		SetRefreshExpired(7200),
	)
}

// ── tests ──

func TestGenerateToken(t *testing.T) {
	store := newFakeStore()
	a := newTestJWTAuth(store)

	tokenInfo, err := a.GenerateToken("user-1", "admin")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	accessToken := tokenInfo.GetAccessToken()
	refreshToken := tokenInfo.GetRefreshToken()

	if accessToken == "" {
		t.Error("access_token is empty")
	}
	if refreshToken == "" {
		t.Error("refresh_token is empty")
	}

	// Both tokens should exist in store
	if found, _ := store.Check(accessToken); !found {
		t.Error("access_token not found in store")
	}
	if found, _ := store.Check(refreshToken); !found {
		t.Error("refresh_token not found in store")
	}
}

func TestParseToken_AccessToken(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		tokenType string
		wantID    string
		wantName  string
		wantErr   error
	}{
		{
			name:      "valid access_token",
			tokenType: "access_token",
			wantErr:   nil,
		},
		{
			name:      "empty token",
			token:     "",
			tokenType: "access_token",
			wantErr:   auth.ErrInvalidToken,
		},
		{
			name:      "destroyed token",
			tokenType: "access_token",
			wantErr:   auth.ErrInvalidToken,
		},
		{
			name:      "unknown token_type",
			tokenType: "unknown",
			wantErr:   auth.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newFakeStore()
			a := newTestJWTAuth(store)

			var token string
			switch tt.name {
			case "valid access_token":
				tokenInfo, err := a.GenerateToken("user-1", "admin")
				if err != nil {
					t.Fatalf("setup: GenerateToken failed: %v", err)
				}
				token = tokenInfo.GetAccessToken()
				tt.wantID = "user-1"
				tt.wantName = "admin"
			case "destroyed token":
				tokenInfo, err := a.GenerateToken("user-2", "test")
				if err != nil {
					t.Fatalf("setup: GenerateToken failed: %v", err)
				}
				token = tokenInfo.GetAccessToken()
				if err := a.DestroyToken(token); err != nil {
					t.Fatalf("setup: DestroyToken failed: %v", err)
				}
			default:
				token = tt.token
			}

			userID, username, err := a.ParseToken(token, tt.tokenType)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if userID != tt.wantID {
				t.Errorf("userID = %q, want %q", userID, tt.wantID)
			}
			if username != tt.wantName {
				t.Errorf("username = %q, want %q", username, tt.wantName)
			}
		})
	}
}

func TestParseToken_RefreshToken(t *testing.T) {
	store := newFakeStore()
	a := newTestJWTAuth(store)

	tokenInfo, err := a.GenerateToken("user-3", "guest")
	if err != nil {
		t.Fatalf("setup: GenerateToken failed: %v", err)
	}

	userID, username, err := a.ParseToken(tokenInfo.GetRefreshToken(), "refresh_token")
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	// refresh_token 通过结构化 claims 携带 user_id/username，不再依赖 Subject 顺序。
	if userID != "user-3" {
		t.Errorf("userID = %q, want %q", userID, "user-3")
	}
	if username != "guest" {
		t.Errorf("username = %q, want %q", username, "guest")
	}
}

// TestParseToken_UsernameWithDot 验证 username 含点号时不会再因 Subject 字符串拼接
// 与 strings.Split 造成身份错乱。结构化 claims 应原样返回 username。
func TestParseToken_UsernameWithDot(t *testing.T) {
	store := newFakeStore()
	a := newTestJWTAuth(store)

	tokenInfo, err := a.GenerateToken("user-dot", "john.doe")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	userID, username, err := a.ParseToken(tokenInfo.GetAccessToken(), "access_token")
	if err != nil {
		t.Fatalf("ParseToken access failed: %v", err)
	}
	if userID != "user-dot" || username != "john.doe" {
		t.Fatalf("access parsed (%q, %q), want (user-dot, john.doe)", userID, username)
	}

	userID, username, err = a.ParseToken(tokenInfo.GetRefreshToken(), "refresh_token")
	if err != nil {
		t.Fatalf("ParseToken refresh failed: %v", err)
	}
	if userID != "user-dot" || username != "john.doe" {
		t.Fatalf("refresh parsed (%q, %q), want (user-dot, john.doe)", userID, username)
	}
}

// TestParseToken_MalformedClaims 验证旧格式/畸形 claims 返回 ErrInvalidToken，且不会 panic。
func TestParseToken_MalformedClaims(t *testing.T) {
	store := newFakeStore()
	a := newTestJWTAuth(store)

	legacy := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: "user.legacy",
	})
	token, err := legacy.SignedString([]byte("test-key"))
	if err != nil {
		t.Fatalf("sign legacy token: %v", err)
	}
	if err := store.Set(token, time.Hour); err != nil {
		t.Fatalf("store legacy token: %v", err)
	}

	if _, _, err := a.ParseToken(token, "access_token"); err != auth.ErrInvalidToken {
		t.Fatalf("ParseToken legacy claims err = %v, want ErrInvalidToken", err)
	}
}

func TestDestroyToken(t *testing.T) {
	store := newFakeStore()
	a := newTestJWTAuth(store)

	tokenInfo, err := a.GenerateToken("user-4", "destroy")
	if err != nil {
		t.Fatalf("setup: GenerateToken failed: %v", err)
	}

	accessToken := tokenInfo.GetAccessToken()
	if err := a.DestroyToken(accessToken); err != nil {
		t.Fatalf("DestroyToken failed: %v", err)
	}

	// After destroy, ParseToken should fail
	_, _, err = a.ParseToken(accessToken, "access_token")
	if err == nil {
		t.Error("expected error after DestroyToken, got nil")
	}
}

func TestRelease(t *testing.T) {
	store := newFakeStore()
	a := newTestJWTAuth(store)

	if err := a.Release(); err != nil {
		t.Errorf("Release failed: %v", err)
	}

	// nil store should not panic
	aNil := New(nil, nil,
		SetSigningMethod(jwt.SigningMethodHS256),
		SetSigningKey([]byte("test-key")),
		SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
			return []byte("test-key"), nil
		}),
	)
	if err := aNil.Release(); err != nil {
		t.Errorf("Release with nil store failed: %v", err)
	}
}

func TestCallStore_NilStore(t *testing.T) {
	a := New(nil, nil,
		SetSigningMethod(jwt.SigningMethodHS256),
		SetSigningKey([]byte("test-key")),
		SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
			return []byte("test-key"), nil
		}),
	)

	// ParseToken with nil store should skip store check and only do JWT validation
	tokenInfo, err := a.GenerateToken("user-5", "nilstore")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	userID, username, err := a.ParseToken(tokenInfo.GetAccessToken(), "access_token")
	if err != nil {
		t.Fatalf("ParseToken with nil store failed: %v", err)
	}
	if userID != "user-5" {
		t.Errorf("userID = %q, want %q", userID, "user-5")
	}
	if username != "nilstore" {
		t.Errorf("username = %q, want %q", username, "nilstore")
	}
}
