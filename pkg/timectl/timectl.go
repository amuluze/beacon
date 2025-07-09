// Package timectl
// Date: 2024/6/25 09:34
// Author: Amu
// Description:
package timectl

import (
	"amprobe/pkg/command"
	"context"
	"fmt"
	"strconv"
	"strings"
)

func GetTime(ctx context.Context) (int64, error) {
	data, err := command.RunCommand(ctx, "date", "+%s")
	if err != nil {
		return 0, nil
	}
	return strconv.ParseInt(strings.TrimSpace(data), 10, 64)
}

func SetTime(ctx context.Context, ts int64) error {
	if ts < 0 {
		return nil
	}
	if _, err := command.RunCommand(ctx, "date", "-u", "-s", fmt.Sprintf("@%d", ts)); err != nil {
		return err
	}
	if _, err := command.RunCommand(ctx, "hwclock", "-w", "-u"); err != nil {
		return err
	}
	return nil
}

func GetTimeZoneList(ctx context.Context) ([]string, error) {
	tzList, err := command.RunCommand(ctx, "timedatectl", "list-timezones")
	if err != nil {
		return []string{}, err
	}
	return strings.Split(strings.TrimSpace(tzList), "\n"), nil
}

func GetTimeZone(ctx context.Context) (string, error) {
	tz, err := command.RunCommand(ctx, "timedatectl", "show", "-p", "Timezone")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Split(strings.TrimSpace(tz), "=")[1]), nil
}

func SetTimeZone(ctx context.Context, tz string) error {
	if _, err := command.RunCommand(ctx, "timedatectl", "set-timezone", tz); err != nil {
		return err
	}
	return nil
}
