// Package main
// Date: 2024/4/23 20:12
// Author: Amu
// Description:
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/amuluze/amvector/service"
	"github.com/takama/daemon"
)

type Service struct {
	configFile string
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
	clearFunc, err := service.Run(s.configFile)
	if err != nil {
		return
	}

	for range interrupt {
		clearFunc()
		return
	}
}

// Stop amvector bootstrap service
func (s *Service) Stop() {
	fmt.Println("Stopping amvector bootstrap service...")
}

func (s *Service) manager() (string, error) {
	usage := "Usage: amvector install | remove | start | stop | status"
	fmt.Printf("os args: %#v\n", os.Args)
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "install":
			return s.daemon.Install()
		case "remove":
			return s.daemon.Remove()
		case "start":
			return s.daemon.Start()
		case "stop":
			return s.daemon.Stop()
		case "status":
			return s.daemon.Status()
		default:
			return usage, nil
		}
	}

	return s.daemon.Run(s)
}
