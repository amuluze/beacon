// Package service
// Date: 2024/3/6 11:09
// Author: Amu
// Description:
package service

import (
	"github.com/amuluze/amprobe/pkg/auth"
	"github.com/amuluze/amprobe/service/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"github.com/google/wire"

	auditAPI "github.com/amuluze/amprobe/service/audit/api"
	authAPI "github.com/amuluze/amprobe/service/auth/api"
	containerAPI "github.com/amuluze/amprobe/service/container/api"
	hostAPI "github.com/amuluze/amprobe/service/host/api"
)

var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

var _ IRouter = (*Router)(nil)

type IRouter interface {
	Register(app *fiber.App) error
	Prefixes() []string
}

type Router struct {
	config *Config
	auth   auth.Auther

	containerAPI *containerAPI.ContainerAPI
	hostAPI      *hostAPI.HostAPI
	authAPI      *authAPI.AuthAPI
	auditAPI     *auditAPI.AuditAPI

	loggerHandler *LoggerHandler
	termHandler   *TermHandler
}

func (a *Router) RegisterAPI(app *fiber.App) {
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
	app.Get("/ws/term", websocket.New(a.termHandler.Handler))

	if a.config.Auth.Enable {
		app.Use(middleware.UserAuthMiddleware(
			a.auth,
			middleware.AllowPathPrefixSkipper("/v1/index/index"),
			middleware.AllowPathPrefixSkipper("/v1/auth/login"),
			middleware.AllowPathPrefixSkipper("/v1/auth/token_update"),
		))
	}

	v1 := app.Group("v1")
	{
		gIndex := v1.Group("index")
		{
			gIndex.Get("/index", func(c *fiber.Ctx) error {
				return c.SendString("hello world")
			})
		}

		gAuth := v1.Group("auth")
		{
			gAuth.Post("/login", a.authAPI.Login).Name("登录")
			gAuth.Post("/logout", a.authAPI.Logout).Name("登出")
			gAuth.Post("/pass_update", a.authAPI.PassUpdate).Name("更新密码")
			gAuth.Post("/token_update", a.authAPI.TokenUpdate).Name("更新 token")
		}

		gContainer := v1.Group("container")
		{
			gContainer.Get("/containers", a.containerAPI.ContainerList).Name("获取容器列表")
			gContainer.Post("/container_start", a.containerAPI.ContainerStart).Name("启动容器")
			gContainer.Post("/container_stop", a.containerAPI.ContainerStop).Name("停止容器")
			gContainer.Post("/container_restart", a.containerAPI.ContainerRestart).Name("重启容器")
			gContainer.Post("/container_remove", a.containerAPI.ContainerRemove).Name("删除容器")
			gContainer.Get("/images", a.containerAPI.ImageList).Name("获取镜像列表")
			gContainer.Post("/image_remove", a.containerAPI.ImageRemove).Name("删除镜像")
			gContainer.Post("/images_prune", a.containerAPI.ImagesPrune).Name("清理虚悬镜像")
			gContainer.Get("/version", a.containerAPI.Version).Name("获取 Docker 版本信息")
			gContainer.Get("/get_docker_registry_mirrors", a.containerAPI.GetDockerRegistryMirrors).Name("获取 Docker 镜像设置")
			gContainer.Post("/set_docker_registry_mirrors", a.containerAPI.SetDockerRegistryMirrors).Name("更新 Docker 镜像设置")
		}

		gHost := v1.Group("host")
		{
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
			gHost.Get("/get_dns_settings", a.hostAPI.GetDNSSettings).Name("获取 DNS 设置")
			gHost.Post("/set_dns_settings", a.hostAPI.SetDNSSettings).Name("更新 DNS 设置")
			gHost.Get("/get_system_time", a.hostAPI.GetSystemTime).Name("获取系统时间")
			gHost.Post("/set_system_time", a.hostAPI.SetSystemTime).Name("更新系统时间")
			gHost.Get("/get_system_timezone", a.hostAPI.GetSystemTimezone).Name("获取系统时区")
			gHost.Post("/set_system_timezone", a.hostAPI.SetSystemTimezone).Name("更新系统时区")
			gHost.Post("/reboot", a.hostAPI.Reboot).Name("重启系统")
			gHost.Post("/shutdown", a.hostAPI.Shutdown).Name("关闭系统")
		}

		gAudit := v1.Group("audit")
		{
			gAudit.Get("/query", a.auditAPI.AuditQuery).Name("获取审计日志")
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
