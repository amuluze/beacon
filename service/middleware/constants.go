// Package middleware
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package middleware

var OperateEvent = map[string]string{
	"/api/v1/auth/login":                  "登录",
	"/api/v1/auth/logout":                 "登出",
	"/api/v1/auth/pass_update":            "更新密码",
	"/api/v1/auth/token_update":           "刷新token",
	"/api/v1/container/container_start":   "启动容器",
	"/api/v1/container/container_stop":    "停止容器",
	"/api/v1/container/container_remove":  "删除容器",
	"/api/v1/container/container_restart": "重启容器",
	"/api/v1/image/image_remove":          "删除镜像",
	"/api/v1/image/images_prune":          "删除虚悬镜像",
	"/ws":                                 "SSH 连接",
}
