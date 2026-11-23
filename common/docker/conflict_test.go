package docker

import (
	"errors"
	"testing"

	"github.com/docker/docker/errdefs"
)

func TestIsNameConflict(t *testing.T) {
	conflict := errdefs.Conflict(errors.New("the container name is already in use"))
	if !IsNameConflict(conflict) {
		t.Fatal("expected docker conflict error to be detected")
	}
	if IsNameConflict(errors.New("boom")) {
		t.Fatal("plain error must not be treated as name conflict")
	}
	if IsNameConflict(nil) {
		t.Fatal("nil error must not be treated as name conflict")
	}
}
