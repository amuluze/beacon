// Package main
// Description: collia Agent CLI 入口。
//
// 子命令：
//   run         前台运行（调试用）
//   daemon      由 systemd 拉起时执行（通常无需手动调用）
//   install     注册系统服务；执行 preflight（探测磁盘/网卡、回写 config、目录/端口预检）
//   uninstall / remove  注销系统服务并清理二进制/配置/数据/日志
//   start/stop/status  通过 daemon 控制 systemd 单元
//   probe       仅执行 preflight 并打印结果（不注册服务）
//   version     输出编译期版本
//
// 设计对齐：spec agent-lifecycle-update.md（RU004 卸载清理）、agent_install.go 的调用契约。
package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"collia/service"

	"github.com/takama/daemon"
)

// version 由 Taskfile 通过 ldflags 注入（-X main.version=...）。
var version = "dev"

const (
	serviceName  = "collia"
	serviceDesc  = "Beacon Probe Agent"
	defaultCfg   = "/etc/collia/config.yml"
	binaryPath   = "/usr/sbin/collia"
	configDir    = "/etc/collia"
	dataDir      = "/data/beacon/collia"
	logDir       = "/data/beacon/logs/collia"
	defaultPrefix = "/"
)

func main() {
	if len(os.Args) < 2 {
		// 无参数：视为被 service manager 拉起
		runAsService(defaultCfg)
		return
	}

	cmd, rest := os.Args[1], os.Args[2:]
	switch cmd {
	case "run":
		cfg := configFromArgs(rest)
		runForeground(cfg)
	case "daemon":
		// systemd ExecStart 目标：与无参数行为一致，显式子命令便于 unit 可读
		cfg := configFromArgs(rest)
		runAsService(cfg)
	case "install":
		cfg := configFromArgs(rest)
		installService(cfg)
	case "uninstall", "remove":
		removeService()
	case "start":
		control("start")
	case "stop":
		control("stop")
	case "status":
		control("status")
	case "probe":
		cfg := configFromArgs(rest)
		probeOnly(cfg)
	case "version", "-v", "--version":
		fmt.Println(version)
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		usage()
		os.Exit(2)
	}
}

// configFromArgs 解析 -c / --config；未命中则用 defaultCfg。
func configFromArgs(args []string) string {
	if cfg := scanConfig(args); cfg != "" {
		return cfg
	}
	return defaultCfg
}

// scanConfig 简单扫描 -c / --config / -config 值；未命中返回空。
func scanConfig(args []string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-c" || a == "--config" || a == "-config":
			if i+1 < len(args) {
				return args[i+1]
			}
		case len(a) > len("--config=") && a[:len("--config=")] == "--config=":
			return a[len("--config="):]
		case len(a) > len("-c=") && a[:len("-c=")] == "-c=":
			return a[len("-c="):]
		}
	}
	return ""
}

func usage() {
	fmt.Fprintf(os.Stderr, `collia (%s) — Beacon Probe Agent

Usage:
  collia <command> [flags]

Commands:
  install                注册系统服务并执行环境预检（探测磁盘/网卡、回写 config）
  uninstall, remove      注销服务并清理二进制/配置/数据/日志
  start | stop | status  控制 systemd 单元
  run                    前台运行（调试）
  daemon                 service manager 拉起入口（通常不手动调用）
  probe                  仅执行预检并打印结果
  version                输出版本号

Flags:
  -c, --config <path>    配置文件路径（默认 %s）
`, version, defaultCfg)
}

// --- 前台/服务运行 ---

func runForeground(cfg string) {
	cleanup, err := service.Run(cfg, defaultPrefix, version)
	if err != nil {
		slog.Error("run failed", "err", err)
		os.Exit(1)
	}
	waitSignal(cleanup)
}

// colliaService 实现 daemon.Executable。
type colliaService struct {
	cfg     string
	cleanup func()
}

func (c *colliaService) Start() {
	cleanup, err := service.Run(c.cfg, defaultPrefix, version)
	if err != nil {
		slog.Error("service run failed", "err", err)
		os.Exit(1)
	}
	c.cleanup = cleanup
}

func (c *colliaService) Stop() {
	if c.cleanup != nil {
		c.cleanup()
	}
}

// Run 由 daemon.Run 调用，负责信号处理与生命周期编排。
func (c *colliaService) Run() {
	c.Start()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
	c.Stop()
}

func runAsService(cfg string) {
	d, err := daemon.New(serviceName, serviceDesc, daemon.SystemDaemon)
	if err != nil {
		slog.Error("create daemon failed", "err", err)
		os.Exit(1)
	}
	svc := &colliaService{cfg: cfg}
	if _, err := d.Run(svc); err != nil {
		slog.Error("daemon run failed", "err", err)
		os.Exit(1)
	}
}

func waitSignal(cleanup func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
	if cleanup != nil {
		cleanup()
	}
}

// --- install / remove / control ---

func installService(cfg string) {
	// 先做环境预检：探测磁盘/网卡、回写 config、目录/端口校验
	if err := preflight(cfg); err != nil {
		// preflight 失败仅警告，不阻断安装（允许在受限环境强行安装）
		slog.Warn("preflight reported issues, continuing install", "err", err)
	}

	d, err := daemon.New(serviceName, serviceDesc, daemon.SystemDaemon)
	if err != nil {
		slog.Error("create daemon failed", "err", err)
		os.Exit(1)
	}
	// systemd 单元 ExecStart 指向 `collia daemon -c <cfg>`
	out, err := d.Install("daemon", "-c", cfg)
	if err != nil {
		// 已安装视为成功（幂等），与 agent_install.go 的 `collia install || true` 一致
		fmt.Fprintln(os.Stderr, out)
		os.Exit(0)
	}
	fmt.Println(out)
}

func removeService() {
	// 先注销 systemd 单元（失败仅警告，单元可能不存在）
	if d, err := daemon.New(serviceName, serviceDesc, daemon.SystemDaemon); err == nil {
		if out, err := d.Remove(); err != nil {
			slog.Warn("daemon remove failed (continuing file cleanup)", "err", err)
		} else {
			fmt.Println(out)
		}
	}
	// 清理文件（spec RU004 / T3-02）
	residuals := cleanupPaths()
	if len(residuals) > 0 {
		fmt.Fprintln(os.Stderr, "部分路径清理失败：")
		for _, r := range residuals {
			fmt.Fprintln(os.Stderr, "  -", r)
		}
		os.Exit(1)
	}
	fmt.Println("collia removed")
}

func control(action string) {
	d, err := daemon.New(serviceName, serviceDesc, daemon.SystemDaemon)
	if err != nil {
		slog.Error("create daemon failed", "err", err)
		os.Exit(1)
	}
	var out string
	switch action {
	case "start":
		out, err = d.Start()
	case "stop":
		out, err = d.Stop()
	case "status":
		out, err = d.Status()
	}
	fmt.Println(out)
	if err != nil {
		os.Exit(1)
	}
}

// --- preflight ---

// preflight 加载 config、探测设备、回写、校验环境。返回聚合 error。
func preflight(cfg string) error {
	return runPreflight(cfg, false)
}

func probeOnly(cfg string) {
	_ = runPreflight(cfg, true)
}

func runPreflight(cfg string, verbose bool) error {
	config, err := service.LoadConfig(cfg)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	report := service.Preflight(config, "")
	if verbose || report.HasError() {
		service.PrintPreflight(report)
	}
	// 合并探测结果（磁盘/网卡）回写 config
	if mErr := service.MergeProbeIntoConfig(cfg, report); mErr != nil {
		slog.Warn("merge probe into config failed", "err", mErr)
	}
	if report.HasError() {
		return fmt.Errorf("preflight 发现阻断性问题：%v", report.Errors)
	}
	return nil
}

// cleanupPaths 删除二进制/配置/数据/日志，返回失败的残留项（spec RU004）。
func cleanupPaths() []string {
	var residuals []string
	for _, p := range []string{binaryPath, binaryPath + ".bak", configDir, dataDir, logDir} {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		if err := os.RemoveAll(p); err != nil {
			residuals = append(residuals, fmt.Sprintf("%s: %v", p, err))
		}
	}
	return residuals
}
