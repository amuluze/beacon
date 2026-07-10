// Package service
// Date: 2024/3/6 11:09
// Author: Amu
// Description:
package service

import (
	"time"

	"beacon/pkg/auth"
	"beacon/pkg/contextx"
	"beacon/service/agent"
	healthapi "beacon/service/health/api"
	"beacon/service/middleware"
	"beacon/service/report"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/wire"

	accountAPI "beacon/service/account/api"
	alarmAPI "beacon/service/alarm/api"
	auditAPI "beacon/service/audit/api"
	authAPI "beacon/service/auth/api"
	containerAPI "beacon/service/container/api"
	hostAPI "beacon/service/host/api"
	mailAPI "beacon/service/mail/api"
)

var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

var _ IRouter = (*Router)(nil)

type IRouter interface {
	Register(app *fiber.App) error
	Prefixes() []string
}

type Router struct {
	config   *Config
	auth     auth.Auther
	enforcer *casbin.SyncedEnforcer

	containerAPI *containerAPI.ContainerAPI
	hostAPI      *hostAPI.HostAPI
	authAPI      *authAPI.AuthAPI
	auditAPI     *auditAPI.AuditAPI
	accountAPI   *accountAPI.AccountAPI
	mailAPI      *mailAPI.MailAPI
	alarmAPI     *alarmAPI.AlarmAPI
	agentAPI     *agent.API
	reportSvc    *report.Service

	loggerHandler *LoggerHandler
	termHandler   *TermHandler

	healthProbe *healthapi.Probe
}

func (a *Router) RegisterAPI(app *fiber.App) {
	a.registerHealthProbes(app)
	a.registerMiddlewares(app)

	api := app.Group("/api")
	a.registerAPIRoutes(api)
}

// registerHealthProbes sets up health check endpoints (no auth required).
func (a *Router) registerHealthProbes(app *fiber.App) {
	app.Get("/health", a.healthProbe.Liveness)
	app.Get("/ready", a.healthProbe.Readiness)
}

// registerMiddlewares sets up global middlewares: rate limiting, agent ID injection,
// WebSocket upgrade, auth, and casbin.
// commonAuthSkipperPaths lists paths that are exempt from authentication and authorization checks.
// These are typically public endpoints (health probes, login, token refresh) and agent installation.
var commonAuthSkipperPaths = []string{"/health", "/ready", "/api/v1/index/index", "/api/v1/auth/login", "/api/v1/auth/token_update", "/api/v1/host/install"}

func (a *Router) registerMiddlewares(app *fiber.App) {
	app.Use("/api/v1/auth/login", limiter.New(limiter.Config{
		Max:        10,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many requests"})
		},
	}))
	app.Use("/api/v1/host/report", limiter.New(limiter.Config{
		Max:        60,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			agentID := c.Get("X-Agent-ID")
			if agentID == "" {
				agentID = c.Query("agent_id")
			}
			if agentID != "" {
				return agentID
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many requests"})
		},
	}))

	app.Use(func(c *fiber.Ctx) error {
		agentID := c.Get("X-Agent-ID")
		if agentID == "" {
			agentID = c.Query("agent_id")
		}
		if agentID != "" {
			c.SetUserContext(contextx.NewAgentID(c.UserContext(), agentID))
		}
		return c.Next()
	})

	app.Use("ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/:id", websocket.New(a.loggerHandler.Handler))
	app.Get("/ws", websocket.New(a.termHandler.Handler))

	if a.config.Auth.Enable {
		app.Use(middleware.UserAuthMiddleware(
			a.auth,
			middleware.AllowPathPrefixSkipper(commonAuthSkipperPaths...),
		))
	}

	if a.config.Casbin.Enable {
		app.Use(middleware.CasbinMiddleware(
			a.enforcer,
			middleware.AllowPathPrefixSkipper(commonAuthSkipperPaths...),
		))
	}
}

// registerAPIRoutes registers all v1 API route groups.
func (a *Router) registerAPIRoutes(api fiber.Router) {
	v1 := api.Group("v1")
	a.registerIndexRoutes(v1)
	a.registerAccountRoutes(v1)
	a.registerAgentListRoute(v1)
	a.registerAuthRoutes(v1)
	a.registerContainerRoutes(v1)
	a.registerHostRoutes(v1)
	a.registerAuditRoutes(v1)
	a.registerMailRoutes(v1)
	a.registerAlarmRoutes(v1)
}

func (a *Router) registerIndexRoutes(v1 fiber.Router) {
	g := v1.Group("index")
	g.Get("/index", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	}).Name("测试")
}

func (a *Router) registerAccountRoutes(v1 fiber.Router) {
	gUser := v1.Group("user")
	gUser.Get("/user_query", a.accountAPI.UserQuery).Name("查询用户")
	gUser.Post("/user_create", a.accountAPI.UserCreate).Name("创建用户")
	gUser.Post("/user_update", a.accountAPI.UserUpdate).Name("更新用户")
	gUser.Post("/user_delete", a.accountAPI.UserDelete).Name("删除用户")

	gRole := v1.Group("role")
	gRole.Get("/role_query", a.accountAPI.RoleQuery).Name("查询角色")
	gRole.Post("/role_create", a.accountAPI.RoleCreate).Name("创建角色")
	gRole.Post("/role_update", a.accountAPI.RoleUpdate).Name("更新角色")
	gRole.Post("/role_delete", a.accountAPI.RoleDelete).Name("删除角色")

	gResource := v1.Group("resource")
	gResource.Get("/resource_query", a.accountAPI.ResourceQuery).Name("查询资源")
}

func (a *Router) registerAgentListRoute(v1 fiber.Router) {
	g := v1.Group("agent")
	g.Get("/list", a.agentAPI.List).Name("查询 Agent 列表")
}

func (a *Router) registerAuthRoutes(v1 fiber.Router) {
	g := v1.Group("auth")
	g.Post("/login", a.authAPI.Login).Name("登录")
	g.Post("/logout", a.authAPI.Logout).Name("登出")
	g.Post("/pass_update", a.authAPI.PassUpdate).Name("更新密码")
	g.Post("/token_update", a.authAPI.TokenUpdate).Name("更新 token")
	g.Get("/user_info", a.authAPI.UserInfo).Name("查询权限")
}

func (a *Router) registerContainerRoutes(v1 fiber.Router) {
	g := v1.Group("container")
	g.Get("/version", a.containerAPI.Version).Name("获取 Docker 版本信息")
	g.Get("/containers", a.containerAPI.ContainerList).Name("获取容器列表")
	g.Get("/usage", a.containerAPI.Usage).Name("获取容器资源使用情况")
	g.Post("/container_create", a.containerAPI.ContainerCreate).Name("创建容器")
	g.Post("/container_update", a.containerAPI.ContainerUpdate).Name("编辑容器")
	g.Post("/container_start", a.containerAPI.ContainerStart).Name("启动容器")
	g.Post("/container_stop", a.containerAPI.ContainerStop).Name("停止容器")
	g.Post("/container_restart", a.containerAPI.ContainerRestart).Name("重启容器")
	g.Post("/container_remove", a.containerAPI.ContainerRemove).Name("删除容器")
	g.Get("/images", a.containerAPI.ImageList).Name("获取镜像列表")
	g.Post("/image_remove", a.containerAPI.ImageRemove).Name("删除镜像")
	g.Post("/images_prune", a.containerAPI.ImagesPrune).Name("清理虚悬镜像")
	g.Post("/image_pull", a.containerAPI.ImagePull).Name("拉取镜像")
	g.Post("/image_import", a.containerAPI.ImageImport).Name("导入镜像")
	g.Post("/image_export", a.containerAPI.ImageExport).Name("导出镜像")
	g.Post("/network_create", a.containerAPI.NetworkCreate).Name("创建网络")
	g.Post("/network_delete", a.containerAPI.NetworkDelete).Name("删除网络")
	g.Get("/networks", a.containerAPI.NetworkList).Name("获取网络列表")
	g.Get("/get_docker_registry_mirrors", a.containerAPI.GetDockerRegistryMirrors).Name("获取 Docker 镜像设置")
	g.Post("/set_docker_registry_mirrors", a.containerAPI.SetDockerRegistryMirrors).Name("更新 Docker 镜像设置")
}

func (a *Router) registerHostRoutes(v1 fiber.Router) {
	g := v1.Group("host")
	g.Get("/install", a.AgentInstallScript).Name("获取 Collia 安装脚本")
	g.Get("/install/package", a.AgentInstallPackage).Name("下载 Collia 安装包")
	g.Get("/install/config", a.AgentInstallConfig).Name("下载 Collia 配置")
	g.Get("/install/certs", a.AgentInstallCerts).Name("下载 Collia 证书")
	g.Post("/get_install_token", a.AgentInstallToken).Name("获取 Collia 安装令牌")
	g.Post("/report", a.reportSvc.HandleReport).Name("Agent 上报监控数据")
	g.Get("/host_info", a.hostAPI.HostInfo).Name("获取主机信息")
	g.Get("/cpu_info", a.hostAPI.CPUInfo).Name("获取 CPU 信息")
	g.Get("/mem_info", a.hostAPI.MemInfo).Name("获取内存信息")
	g.Get("/disk_info", a.hostAPI.DiskInfo).Name("获取磁盘信息")
	g.Get("/cpu_trending", a.hostAPI.CPUUsage).Name("获取 CPU 使用率")
	g.Get("/mem_trending", a.hostAPI.MemUsage).Name("获取内存使用率")
	g.Get("/disk_trending", a.hostAPI.DiskUsage).Name("获取磁盘使用率")
	g.Get("/net_trending", a.hostAPI.NetUsage).Name("获取网络使用率")
	g.Get("/file_search", a.hostAPI.FilesSearch).Name("文件搜索")
	g.Post("/file_upload", a.hostAPI.FileUpload).Name("文件上传")
	g.Post("/file_download", a.hostAPI.FileDownload).Name("文件下载")
	g.Post("/file_delete", a.hostAPI.FileDelete).Name("删除文件")
	g.Post("/file_create", a.hostAPI.FileCreate).Name("创建文件")
	g.Post("/folder_create", a.hostAPI.FolderCreate).Name("创建文件夹")
	g.Get("/get_dns_settings", a.hostAPI.GetDNSSettings).Name("获取 DNS 设置")
	g.Post("/set_dns_settings", a.hostAPI.SetDNSSettings).Name("更新 DNS 设置")
	g.Get("/get_system_time", a.hostAPI.GetSystemTime).Name("获取系统时间")
	g.Post("/set_system_time", a.hostAPI.SetSystemTime).Name("更新系统时间")
	g.Get("/get_system_timezone_list", a.hostAPI.GetSystemTimeZoneList).Name("获取系统时区列表")
	g.Get("/get_system_timezone", a.hostAPI.GetSystemTimeZone).Name("获取系统时区")
	g.Post("/set_system_timezone", a.hostAPI.SetSystemTimeZone).Name("更新系统时区")
	g.Post("/reboot", a.hostAPI.Reboot).Name("重启系统")
	g.Post("/shutdown", a.hostAPI.Shutdown).Name("关闭系统")
}

func (a *Router) registerAuditRoutes(v1 fiber.Router) {
	g := v1.Group("audit")
	g.Get("/query", a.auditAPI.AuditQuery).Name("获取审计日志")
}

func (a *Router) registerMailRoutes(v1 fiber.Router) {
	g := v1.Group("mail")
	g.Post("/mail_create", a.mailAPI.MailCreate).Name("创建邮件告警配置")
	g.Post("/mail_delete", a.mailAPI.MailDelete).Name("删除邮件告警配置")
	g.Post("/mail_update", a.mailAPI.MailUpdate).Name("更新邮件告警配置")
	g.Get("/mail_query", a.mailAPI.MailQuery).Name("查询邮件告警配置")
	g.Post("/mail_test", a.mailAPI.MailTest).Name("测试邮件告警")
}

func (a *Router) registerAlarmRoutes(v1 fiber.Router) {
	g := v1.Group("alarm")
	g.Post("/alarm_update", a.alarmAPI.AlarmUpdate).Name("更新告警阈值")
	g.Get("/alarm_query", a.alarmAPI.AlarmQuery).Name("查询告警阈值")
}

func (a *Router) Register(app *fiber.App) error {
	a.RegisterAPI(app)
	return nil
}
func (a *Router) Prefixes() []string {
	return []string{"/api/"}
}
