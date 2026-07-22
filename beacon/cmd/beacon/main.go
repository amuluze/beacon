// Package main
// Date: 2024/3/6 11:21
// Author: Amu
// Description:
package main

import (
	"context"
	"os"

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
