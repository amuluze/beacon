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

	"collia/service"

	"github.com/takama/daemon"
)

const (
	name        = "collia"
	description = "beacon agent service"
)

var dependencies = []string{""}

var (
	configFile string
	prefix     string
	Version    string // set via -ldflags at compile time
)

func usage() {
	fmt.Println("Description: \n\t", description)
	fmt.Println("Usage: \n\t", os.Args[0], " [--flag arguments] install | remove | start | stop | status")
	fmt.Println("Flags: ")
	flag.PrintDefaults()
}

func parseConfig() []string {
	flag.StringVar(&configFile, "conf", "/etc/collia/config.yml", "config file path")
	flag.StringVar(&prefix, "prefix", "/data/beacon", "prefix of beacon server-agent resources dir")
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
	svc := &Service{daemon: src, configFile: configFile, prefix: service.Prefix(prefix), version: Version}
	status, err := svc.manager(args)
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
