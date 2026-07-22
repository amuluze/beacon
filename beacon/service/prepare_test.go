package service

import (
	"path/filepath"
	"testing"

	"beacon/pkg/utils/hash"
	"beacon/service/model"
	"common/database"
)

func TestInitAlarmThresholdCreatesAllDefaultsAndPreservesOverrides(t *testing.T) {
	db, err := database.NewDB(
		database.WithType("sqlite"),
		database.WithDBName(filepath.Join(t.TempDir(), "beacon")),
	)
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(&model.AlarmThreshold{}); err != nil {
		t.Fatalf("migrate alarm threshold: %v", err)
	}

	prepare := &Prepare{db: db}
	prepare.InitAlarmThreshold()

	var got []model.AlarmThreshold
	if err := db.Order("type").Find(&got).Error; err != nil {
		t.Fatalf("query alarm thresholds: %v", err)
	}
	if len(got) != len(thresholds) {
		t.Fatalf("alarm threshold count = %d, want %d", len(got), len(thresholds))
	}

	if err := db.Model(&model.AlarmThreshold{}).
		Where("type = ?", "memory").
		Updates(map[string]any{"duration": 5, "threshold": 90}).Error; err != nil {
		t.Fatalf("override memory alarm threshold: %v", err)
	}
	prepare.InitAlarmThreshold()

	var count int64
	if err := db.Model(&model.AlarmThreshold{}).Count(&count).Error; err != nil {
		t.Fatalf("count alarm thresholds: %v", err)
	}
	if count != int64(len(thresholds)) {
		t.Fatalf("alarm threshold count after reinitialization = %d, want %d", count, len(thresholds))
	}

	var memory model.AlarmThreshold
	if err := db.Where("type = ?", "memory").First(&memory).Error; err != nil {
		t.Fatalf("query memory alarm threshold: %v", err)
	}
	if memory.Duration != 5 || memory.Threshold != 90 {
		t.Fatalf("memory override was replaced: duration=%d threshold=%d", memory.Duration, memory.Threshold)
	}
}

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
