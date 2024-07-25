// Package middleware
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package middleware

var OperateEvent = map[string]string{
	"/v1/auth/login":                            "登录",
	"/v1/auth/logout":                           "登出",
	"/v1/auth/pass_update":                      "更新密码",
	"/v1/auth/token_update":                     "刷新token",
	"/v1/container/container_create":            "创建容器",
	"/v1/container/container_start":             "启动容器",
	"/v1/container/container_stop":              "停止容器",
	"/v1/container/container_remove":            "删除容器",
	"/v1/container/container_restart":           "重启容器",
	"/v1/container/image_remove":                "删除镜像",
	"/v1/container/images_prune":                "删除虚悬镜像",
	"/v1/container/images_pull":                 "拉取镜像",
	"/v1/container/images_import":               "导入镜像",
	"/v1/container/images_export":               "导出镜像",
	"/v1/container/network_create":              "创建网络",
	"/v1/container/network_delete":              "删除网络",
	"/v1/container/set_docker_registry_mirrors": "更新docker镜像加速设置",
	"/v1/host/file_upload":                      "上传文件",
	"/v1/host/file_download":                    "下载文件",
	"/v1/host/file_delete":                      "删除文件",
	"/v1/host/file_create":                      "创建文件",
	"/v1/host/folder_create":                    "创建文件夹",
	"/v1/host/set_dns_settings":                 "更新 DNS 设置",
	"/v1/host/set_system_time":                  "更新系统时间",
	"/v1/host/set_system_timezone":              "更新系统时区",
	"/v1/host/reboot":                           "重启",
	"/v1/host/shutdown":                         "关机",
}
