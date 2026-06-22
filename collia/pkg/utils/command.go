// Package utils
// Date: 2024/6/25 10:09
// Author: Amu
// Description:
package utils

import (
	"bytes"
	"context"
	"os/exec"
)

func RunCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	cmd.Stderr = buf
	if err := cmd.Start(); err != nil {
		return buf.String(), err
	}
	if err := cmd.Wait(); err != nil {
		return buf.String(), err
	}
	return buf.String(), nil
}
