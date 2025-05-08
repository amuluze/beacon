// Package main
// Date: 2024/3/6 11:21
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
	name        = "amprobe"
	description = "data collection and monitor services"
)

var dependencies = []string{""}

var (
	configFile string
	modelFile  string
)

func usage() {
	fmt.Println("Description: \n\t", description)
	fmt.Println("Usage: \n\t", os.Args[0], " [--flag arguments] setup | install | remove | start | stop | status")
	fmt.Println("Flags: ")
	flag.PrintDefaults()
}

func parseConfig() []string {
	flag.StringVar(&configFile, "conf", "/etc/amprobe/config.yml", "config file path")
	flag.StringVar(&modelFile, "model", "/etc/amprobe/model.conf", "model file path")
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
	srv := &Service{daemon: src, configFile: configFile, modelFile: modelFile}
	status, err := srv.manager(args)
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
