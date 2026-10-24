// Package docker
// Date: 2024/07/09 14:14:04
// Author: Amu
// Description:
package docker

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestImageLifecycle(t *testing.T) {
	manager := testManager(t)
	imageName := ensureTestImage(t, manager)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	images, err := manager.ListImage(ctx)
	if err != nil {
		t.Fatalf("list images: %v", err)
	}
	if len(images) == 0 {
		t.Fatal("Docker daemon returned no tagged images")
	}
	summary, err := manager.GetImageByName(ctx, imageName)
	if err != nil {
		t.Fatalf("get image by name: %v", err)
	}
	if _, err := manager.GetImageByID(ctx, summary.ID); err != nil {
		t.Fatalf("get image by ID: %v", err)
	}

	tag := uniqueDockerName("image") + ":test"
	if err := manager.TagImage(ctx, imageName, tag); err != nil {
		t.Fatalf("tag image: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cleanupCancel()
		if err := manager.DeleteImage(cleanupCtx, tag); err != nil {
			t.Logf("cleanup image tag %s: %v", tag, err)
		}
	})
	if _, err := manager.GetImageByName(ctx, tag); err != nil {
		t.Fatalf("get tagged image: %v", err)
	}

	targetFile := filepath.Join(t.TempDir(), "image.tar")
	if err := os.WriteFile(targetFile, []byte("stale image"), 0o644); err != nil {
		t.Fatalf("seed image archive: %v", err)
	}
	if err := os.Chmod(targetFile, 0o644); err != nil {
		t.Fatalf("set insecure image archive permissions: %v", err)
	}
	if err := manager.ExportImage(ctx, []string{tag}, targetFile); err != nil {
		t.Fatalf("export image: %v", err)
	}
	info, err := os.Stat(targetFile)
	if err != nil {
		t.Fatalf("stat exported image: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("exported image permissions = %o, want 600", got)
	}
	if err := manager.ImportImage(ctx, targetFile); err != nil {
		t.Fatalf("import image: %v", err)
	}
}
