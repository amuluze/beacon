// Package schema
// Date: 2025/02/12 15:18:59
// Author: Amu
// Description:
package schema

type Statistics struct {
	ID    uint `json:"id"`
	Times int  `json:"times"`
}

type StatisticsQueryArgs struct {
}

type StatisticsQueryReply struct {
	Data Statistics `json:"data"`
}

// StatisticsUpdateArgs 仅允许对存在的统计记录自增；ID 为 0 视为非法。
type StatisticsUpdateArgs struct {
	ID uint `json:"id" validate:"required"`
}

type StatisticsUpdateReply struct {
}

// InstallationReportArgs 来自安装脚本的一次性上报，所有字符串字段限长以防滥用。
type InstallationReportArgs struct {
	InstallID     string `json:"install_id" validate:"required,max=128"`
	Image         string `json:"image" validate:"max=256"`
	Version       string `json:"version" validate:"max=64"`
	PublicBaseURL string `json:"public_base_url" validate:"max=512"`
	InstallDir    string `json:"install_dir" validate:"max=512"`
	HTTPPort      string `json:"http_port" validate:"max=16"`
	ControlPort   string `json:"control_port" validate:"max=16"`
	ContainerName string `json:"container_name" validate:"max=128"`
	Hostname      string `json:"hostname" validate:"max=256"`
	ClientIP      string `json:"-"`
	UserAgent     string `json:"-"`
}

type InstallationReportReply struct {
}
