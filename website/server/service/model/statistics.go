// Package model
// Date: 2025/02/12 15:06:51
// Author: Amu
// Description:
package model

import "gorm.io/gorm"

type Statistics struct {
	gorm.Model
	Times int `gorm:"comment:下载次数"`
}

func (d *Statistics) TableName() string {
	return "s_statistics"
}

type InstallationReport struct {
	gorm.Model
	InstallID     string `gorm:"uniqueIndex;size:128;comment:安装标识"`
	Image         string `gorm:"size:255;comment:镜像"`
	Version       string `gorm:"size:64;comment:版本"`
	PublicBaseURL string `gorm:"size:255;comment:公开访问地址"`
	InstallDir    string `gorm:"size:255;comment:安装目录"`
	HTTPPort      string `gorm:"size:32;comment:Web端口"`
	ControlPort   string `gorm:"size:32;comment:控制端口"`
	ContainerName string `gorm:"size:128;comment:容器名称"`
	Hostname      string `gorm:"size:128;comment:容器主机名"`
	ClientIP      string `gorm:"size:64;comment:客户端IP"`
	UserAgent     string `gorm:"size:255;comment:User-Agent"`
}

func (d *InstallationReport) TableName() string {
	return "s_installation_report"
}
