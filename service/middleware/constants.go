// Package middleware
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package middleware

var OperateEvent = map[string]string{
	"/v1/auth/login":                  "登录",
	"/v1/auth/logout":                 "登出",
	"/v1/auth/pass_update":            "更新密码",
	"/v1/auth/token_update":           "刷新token",
	"/v1/container/container_start":   "启动容器",
	"/v1/container/container_stop":    "停止容器",
	"/v1/container/container_remove":  "删除容器",
	"/v1/container/container_restart": "重启容器",
	"/v1/container/image_remove":      "删除镜像",
	"/v1/container/images_prune":      "删除虚悬镜像",
}
