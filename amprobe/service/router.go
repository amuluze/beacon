// Package service
// Date: 2024/3/6 11:09
// Author: Amu
// Description:
package service

import (
	"time"

	"amprobe/pkg/auth"
	"amprobe/pkg/contextx"
	"amprobe/service/agent"
	"amprobe/service/health"
	"amprobe/service/middleware"
	"amprobe/service/report"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/wire"

	accountAPI "amprobe/service/account/api"
	auditAPI "amprobe/service/audit/api"
	authAPI "amprobe/service/auth/api"
	containerAPI "amprobe/service/container/api"
	hostAPI "amprobe/service/host/api"
	mailAPI "amprobe/service/mail/api"
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

	healthProbe *health.Probe
}

func (a *Router) RegisterAPI(app *fiber.App) {
	// Health probes — no auth required, must be placed before auth middleware.
	if a.healthProbe == nil {
		a.healthProbe = health.NewProbe()
	}
	app.Get("/health", a.healthProbe.Liveness)
	app.Get("/ready", a.healthProbe.Readiness)

	// Rate limiting for authentication and Agent reporting endpoints.
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
		// Do NOT fall back to DefaultAgentID; let repository RequireAgentID
		// reject requests that truly need an agent scope with a 400 Bad Request.
		if agentID != "" {
			c.SetUserContext(contextx.NewAgentID(c.UserContext(), agentID))
		}
		return c.Next()
	})

	app.Use("ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
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
			middleware.AllowPathPrefixSkipper("/api/v1/index/index"),
			middleware.AllowPathPrefixSkipper("/api/v1/auth/login"),
			middleware.AllowPathPrefixSkipper("/api/v1/auth/token_update"),
			middleware.AllowPathPrefixSkipper("/api/v1/host/install"),
		))
	}

	if a.config.Casbin.Enable {
		app.Use(middleware.CasbinMiddleware(
			a.enforcer,
			middleware.AllowPathPrefixSkipper("/api/v1/index/index"),
			middleware.AllowPathPrefixSkipper("/api/v1/auth/login"),
			middleware.AllowPathPrefixSkipper("/api/v1/auth/token_update"),
			middleware.AllowPathPrefixSkipper("/api/v1/host/install"),
		))
	}
	api := app.Group("/api")
	{
		v1 := api.Group("v1")
		{
			gIndex := v1.Group("index")
			{
				gIndex.Get("/index", func(c *fiber.Ctx) error {
					return c.SendString("hello world")
				}).Name("测试")
			}

			gUser := v1.Group("user")
			{
				gUser.Get("/user_query", a.accountAPI.UserQuery).Name("查询用户")
				gUser.Post("/user_create", a.accountAPI.UserCreate).Name("创建用户")
				gUser.Post("/user_update", a.accountAPI.UserUpdate).Name("更新用户")
				gUser.Post("/user_delete", a.accountAPI.UserDelete).Name("删除用户")
			}

			gRole := v1.Group("role")
			{
				gRole.Get("/role_query", a.accountAPI.RoleQuery).Name("查询角色")
				gRole.Post("/role_create", a.accountAPI.RoleCreate).Name("创建角色")
				gRole.Post("/role_update", a.accountAPI.RoleUpdate).Name("更新角色")
				gRole.Post("/role_delete", a.accountAPI.RoleDelete).Name("删除角色")
			}

			gResource := v1.Group("resource")
			{
				gResource.Get("/resource_query", a.accountAPI.ResourceQuery).Name("查询资源")
			}

			gAgent := v1.Group("agent")
			{
				gAgent.Get("/list", a.agentAPI.List).Name("查询 Agent 列表")
			}

			gAuth := v1.Group("auth")
			{
				gAuth.Post("/login", a.authAPI.Login).Name("登录")
				gAuth.Post("/logout", a.authAPI.Logout).Name("登出")
				gAuth.Post("/pass_update", a.authAPI.PassUpdate).Name("更新密码")
				gAuth.Post("/token_update", a.authAPI.TokenUpdate).Name("更新 token")
				gAuth.Get("/user_info", a.authAPI.UserInfo).Name("查询权限")
			}

			gContainer := v1.Group("container")
			{
				gContainer.Get("/version", a.containerAPI.Version).Name("获取 Docker 版本信息")
				gContainer.Get("/containers", a.containerAPI.ContainerList).Name("获取容器列表")
				gContainer.Get("/usage", a.containerAPI.Usage).Name("获取容器资源使用情况")
				gContainer.Post("/container_create", a.containerAPI.ContainerCreate).Name("创建容器")
				gContainer.Post("/container_start", a.containerAPI.ContainerStart).Name("启动容器")
				gContainer.Post("/container_stop", a.containerAPI.ContainerStop).Name("停止容器")
				gContainer.Post("/container_restart", a.containerAPI.ContainerRestart).Name("重启容器")
				gContainer.Post("/container_remove", a.containerAPI.ContainerRemove).Name("删除容器")
				gContainer.Get("/images", a.containerAPI.ImageList).Name("获取镜像列表")
				gContainer.Post("/image_remove", a.containerAPI.ImageRemove).Name("删除镜像")
				gContainer.Post("/images_prune", a.containerAPI.ImagesPrune).Name("清理虚悬镜像")
				gContainer.Post("/image_pull", a.containerAPI.ImagePull).Name("拉取镜像")
				gContainer.Post("/image_import", a.containerAPI.ImageImport).Name("导入镜像")
				gContainer.Post("/image_export", a.containerAPI.ImageExport).Name("导出镜像")
				gContainer.Post("/network_create", a.containerAPI.NetworkCreate).Name("创建网络")
				gContainer.Post("/network_delete", a.containerAPI.NetworkDelete).Name("删除网络")
				gContainer.Get("/networks", a.containerAPI.NetworkList).Name("获取网络列表")
				gContainer.Get("/get_docker_registry_mirrors", a.containerAPI.GetDockerRegistryMirrors).Name("获取 Docker 镜像设置")
				gContainer.Post("/set_docker_registry_mirrors", a.containerAPI.SetDockerRegistryMirrors).Name("更新 Docker 镜像设置")
			}

			gHost := v1.Group("host")
			{
				gHost.Get("/install", a.AgentInstallScript).Name("获取 Collia 安装脚本")
				gHost.Get("/install/package", a.AgentInstallPackage).Name("下载 Collia 安装包")
				gHost.Get("/install/config", a.AgentInstallConfig).Name("下载 Collia 配置")
				gHost.Get("/install/certs", a.AgentInstallCerts).Name("下载 Collia 证书")
				gHost.Post("/report", a.reportSvc.HandleReport).Name("Agent 上报监控数据")
				gHost.Get("/host_info", a.hostAPI.HostInfo).Name("获取主机信息")
				gHost.Get("/cpu_info", a.hostAPI.CPUInfo).Name("获取 CPU 信息")
				gHost.Get("/mem_info", a.hostAPI.MemInfo).Name("获取内存信息")
				gHost.Get("/disk_info", a.hostAPI.DiskInfo).Name("获取磁盘信息")
				gHost.Get("/cpu_trending", a.hostAPI.CPUUsage).Name("获取 CPU 使用率")
				gHost.Get("/mem_trending", a.hostAPI.MemUsage).Name("获取内存使用率")
				gHost.Get("/disk_trending", a.hostAPI.DiskUsage).Name("获取磁盘使用率")
				gHost.Get("/net_trending", a.hostAPI.NetUsage).Name("获取网络使用率")
				gHost.Get("/file_search", a.hostAPI.FilesSearch).Name("文件搜索")
				gHost.Post("/file_upload", a.hostAPI.FileUpload).Name("文件上传")
				gHost.Post("/file_download", a.hostAPI.FileDownload).Name("文件下载")
				gHost.Post("/file_delete", a.hostAPI.FileDelete).Name("删除文件")
				gHost.Post("/file_create", a.hostAPI.FileCreate).Name("创建文件")
				gHost.Post("/folder_create", a.hostAPI.FolderCreate).Name("创建文件夹")
				gHost.Get("/get_dns_settings", a.hostAPI.GetDNSSettings).Name("获取 DNS 设置")
				gHost.Post("/set_dns_settings", a.hostAPI.SetDNSSettings).Name("更新 DNS 设置")
				gHost.Get("/get_system_time", a.hostAPI.GetSystemTime).Name("获取系统时间")
				gHost.Post("/set_system_time", a.hostAPI.SetSystemTime).Name("更新系统时间")
				gHost.Get("/get_system_timezone_list", a.hostAPI.GetSystemTimeZoneList).Name("获取系统时区列表")
				gHost.Get("/get_system_timezone", a.hostAPI.GetSystemTimeZone).Name("获取系统时区")
				gHost.Post("/set_system_timezone", a.hostAPI.SetSystemTimeZone).Name("更新系统时区")
				gHost.Post("/reboot", a.hostAPI.Reboot).Name("重启系统")
				gHost.Post("/shutdown", a.hostAPI.Shutdown).Name("关闭系统")
			}

			gAudit := v1.Group("audit")
			{
				gAudit.Get("/query", a.auditAPI.AuditQuery).Name("获取审计日志")
			}

			gMail := v1.Group("mail")
			{
				gMail.Post("/mail_create", a.mailAPI.MailCreate).Name("创建邮件告警配置")
				gMail.Post("/mail_delete", a.mailAPI.MailDelete).Name("删除邮件告警配置")
				gMail.Post("/mail_update", a.mailAPI.MailUpdate).Name("更新邮件告警配置")
				gMail.Get("/mail_query", a.mailAPI.MailQuery).Name("查询邮件告警配置")
				gMail.Post("/mail_test", a.mailAPI.MailTest).Name("测试邮件告警")
			}

			gAlarm := v1.Group("alarm")
			{
				gAlarm.Post("/alarm_update", a.alarmAPI.AlarmUpdate).Name("更新告警阈值")
				gAlarm.Get("/alarm_query", a.alarmAPI.AlarmQuery).Name("查询告警阈值")
			}
		}
	}
}

func (a *Router) Register(app *fiber.App) error {
	a.RegisterAPI(app)
	return nil
}
func (a *Router) Prefixes() []string {
	return []string{"/api/"}
}
