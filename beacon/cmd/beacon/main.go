// Package main
// Date: 2024/3/6 11:21
// Author: Amu
// Description:
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"beacon/service"

	"github.com/urfave/cli/v2"
)

// Build metadata injected via -ldflags at compile time.
var (
	Version    string
	BuildStamp string
	GitHash    string
	GitBranch  string
)

func main() {
	ctx := context.Background()
	// 把编译期版本同步给 VersionChecker，使 /api/v1/system/update 上报真实版本。
	if Version != "" {
		service.BuildVersion = Version
	}

	app := cli.NewApp()
	if Version != "" {
		app.Version = Version
	} else {
		app.Version = "dev"
	}
	app.Usage = "resource monitor"
	app.Commands = []*cli.Command{
		monitorCmd(ctx),
		resetAdminCmd(),
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func monitorCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "run beacon service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "App Configuration file(.toml)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "model",
				Aliases:  []string{"m"},
				Usage:    "Model file(.conf)",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			return service.Run(
				ctx,
				service.SetConfigFile(c.String("conf")),
				service.SetModelFile(c.String("model")),
			)
		},
	}
}

// resetAdminCmd 重置 admin 密码：默认生成临时密码并强制首登修改；
// 可用 --password 或 --password-stdin 指定明文（须满足密码策略）。
func resetAdminCmd() *cli.Command {
	return &cli.Command{
		Name:  "reset-admin",
		Usage: "reset the administrator password to a generated temporary value (or --password/--password-stdin)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Value:   "configs/config.toml",
				Usage:   "App Configuration file(.toml)",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "explicit new password (at least 12 chars, at least three character classes)",
			},
			&cli.BoolFlag{
				Name:  "password-stdin",
				Usage: "read the new password from stdin (takes precedence over --password)",
			},
		},
		Action: func(c *cli.Context) error {
			password := c.String("password")
			if c.Bool("password-stdin") {
				raw, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("read password from stdin: %w", err)
				}
				password = strings.TrimRight(string(raw), "\r\n")
				if password == "" {
					return errors.New("empty password from stdin")
				}
			}
			password, generated, err := service.ResetAdminFromConfig(c.String("conf"), password)
			if err != nil {
				return err
			}
			if generated {
				fmt.Printf("administrator password reset to a temporary value (change it at first login):\n%s\n", password)
			} else {
				fmt.Println("administrator password reset (change it at first login)")
			}
			return nil
		},
	}
}
