// Package timectl
// Date: 2024/6/25 09:34
// Author: Amu
// Description:
package timectl

import (
	"context"
	"fmt"
	"github.com/amuluze/amvector/pkg/utils"
	"strconv"
	"strings"
)

func GetTime(ctx context.Context) (int64, error) {
	data, err := utils.RunCommand(ctx, "date", "+%s")
	if err != nil {
		return 0, nil
	}
	return strconv.ParseInt(strings.TrimSpace(data), 10, 64)
}

func SetTime(ctx context.Context, ts int64) error {
	if ts < 0 {
		return nil
	}
	if _, err := utils.RunCommand(ctx, "date", "-u", "-s", fmt.Sprintf("@%d", ts)); err != nil {
		return err
	}
	if _, err := utils.RunCommand(ctx, "hwclock", "-w", "-u"); err != nil {
		return err
	}
	return nil
}

func GetTimeZone(ctx context.Context) (string, error) {
	tz, err := utils.RunCommand(ctx, "timedatectl", "show", "-p", "Timezone")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Split(strings.TrimSpace(tz), "=")[1]), nil
}

func SetTimeZone(ctx context.Context, tz string) error {
	if _, err := utils.RunCommand(ctx, "timedatectl", "set-timezone", tz); err != nil {
		return err
	}
	return nil
}
