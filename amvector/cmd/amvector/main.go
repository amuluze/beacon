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

func usage() {
	fmt.Println("Description: \n\t", description)
	fmt.Println("Usage: \n\t", os.Args[0], " [--flag arguments] install | remove | start | stop | status")
	fmt.Println("Flags: ")
	flag.PrintDefaults()
}

func parseConfig() []string {
	flag.StringVar(&configFile, "conf", "/data/amvector/configs/config.toml", "config file path")
	flag.Parse()
	return flag.Args()
}

func main() {
	flag.Usage = usage
	args := parseConfig()

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
	status, err := service.manager(args)
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
