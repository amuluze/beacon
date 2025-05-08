// Package main
// Date       : 2025/5/7 14:53
// Author     : Amu
// Description:
package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"amprobe/service"

	"github.com/takama/daemon"
)

type Service struct {
	configFile string
	modelFile  string
	daemon     daemon.Daemon
}

// Start amvector bootstrap service non-blocking
func (s *Service) Start() {
	fmt.Printf("Starting amvector bootstrap service...\n")
}

// Run start amvector bootstrap service blocking and wait for exit signal
func (s *Service) Run() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	clearFunc, err := service.Run(s.configFile, service.ModelConfig(s.modelFile))
	if err != nil {
		slog.Error("service run error", "error", err)
		return
	}

	for range interrupt {
		slog.Info("received exit signal")
		clearFunc()
		return
	}
}

// Stop amvector bootstrap service
func (s *Service) Stop() {
	fmt.Println("Stopping amvector bootstrap service...")
}

func (s *Service) manager(args []string) (string, error) {
	if len(args) > 0 {
		command := args[0]
		switch command {
		case "install":
			installArgs := args[1:]
			return s.daemon.Install(installArgs...)
		case "remove":
			return s.daemon.Remove()
		case "start":
			return s.daemon.Start()
		case "stop":
			return s.daemon.Stop()
		case "status":
			return s.daemon.Status()
		case "setup":
			fmt.Printf("initializing configuration files...\n")
			if err := runSetup(); err != nil {
				fmt.Printf("error initializing configuration files: %v\n", err)
				os.Exit(-1)
			} else {
				os.Exit(0)
			}
		default:
			usage()
			return "", nil
		}
	}

	return s.daemon.Run(s)
}
