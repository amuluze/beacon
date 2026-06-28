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

type StatisticsUpdateArgs struct {
	ID uint `json:"id"`
}

type StatisticsUpdateReply struct {
}

type InstallationReportArgs struct {
	InstallID     string `json:"install_id"`
	Image         string `json:"image"`
	Version       string `json:"version"`
	PublicBaseURL string `json:"public_base_url"`
	InstallDir    string `json:"install_dir"`
	HTTPPort      string `json:"http_port"`
	ControlPort   string `json:"control_port"`
	ContainerName string `json:"container_name"`
	Hostname      string `json:"hostname"`
	ClientIP      string `json:"-"`
	UserAgent     string `json:"-"`
}

type InstallationReportReply struct {
}
