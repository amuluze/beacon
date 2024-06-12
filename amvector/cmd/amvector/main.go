// Package main
// Date: 2024/4/23 19:33
// Author: Amu
// Description:
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/takama/daemon"
)

const (
	name        = "amvector"
	description = "data collection and synchronization services"
)

var dependencies = []string{""}

var configFile string

func main() {
	flag.StringVar(&configFile, "conf", "/data/amvector/configs/config.toml", "config file path")
	flag.Parse()

	var kind daemon.Kind
	if runtime.GOOS == "darwin" {
		kind = daemon.GlobalDaemon
	} else if runtime.GOOS == "linux" {
		kind = daemon.SystemDaemon
	}

	src, err := daemon.New(name, description, kind, dependencies...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	service := &Service{daemon: src, configFile: configFile}
	status, err := service.manager()
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
