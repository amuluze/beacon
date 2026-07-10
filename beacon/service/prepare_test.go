package service

import (
	"testing"

	"beacon/pkg/utils/hash"
)

func TestDefaultUserPasswordsUseLoginHash(t *testing.T) {
	cases := []struct {
		name     string
		username string
		password string
	}{
		{name: "admin", username: "admin", password: "admin123"},
		{name: "beacon", username: "beacon", password: "123456"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, user := range users {
				if user.Username != tc.username {
					continue
				}
				if err := hash.BcryptVerify(tc.password, user.Password); err != nil {
					t.Fatalf("default password for %q should pass bcrypt login verification: %v", tc.username, err)
				}
				return
			}
			t.Fatalf("default user %q not found", tc.username)
		})
	}
}

func TestLegacyDefaultPasswordHashesMatchPreviousDefaults(t *testing.T) {
	cases := []struct {
		name     string
		username string
		password string
	}{
		{name: "admin", username: "admin", password: "admin123"},
		{name: "beacon", username: "beacon", password: "123456"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := legacyDefaultPasswordHashes[tc.username]
			if !ok {
				t.Fatalf("legacy hash for %q not found", tc.username)
			}
			if want := hash.SHA1String(tc.password); got != want {
				t.Fatalf("legacy hash for %q = %q, want %q", tc.username, got, want)
			}
		})
	}
}
