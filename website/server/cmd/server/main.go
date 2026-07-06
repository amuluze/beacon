// Package main
// Date: 2025/02/12 14:56:57
// Author: Amu
// Description:
package main

import (
	"context"
	"os"
	"server/service"

	"github.com/urfave/cli/v2"
)

func main() {
	ctx := context.TODO()
	app := cli.NewApp()
	app.Version = "0.0.1"
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
		Usage: "run server service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "App Configuration file(.toml)",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			return service.Run(
				ctx,
				service.SetConfigFile(c.String("conf")),
			)
		},
	}
}
