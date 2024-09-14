// Package schema
// Date: 2024/3/6 13:20
// Author: Amu
// Description:
package schema

type Usage struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type NetIO struct {
	Timestamp int64   `json:"timestamp"`
	BytesSent float64 `json:"bytes_sent"`
	BytesRecv float64 `json:"bytes_recv"`
}

type DiskIO struct {
	Timestamp int64   `json:"timestamp"`
	IORead    float64 `json:"io_read"`
	IOWrite   float64 `json:"io_write"`
}

type HostArgs struct{}

type HostInfoReply struct {
	Timestamp       int64  `json:"timestamp"`
	Uptime          string `json:"uptime"`
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	KernelVersion   string `json:"kernel_version"`
	KernelArch      string `json:"kernel_arch"`
}

type CPUArgs struct{}

type CPUInfoReply struct {
	Percent float64 `json:"percent"`
}

type CPUUsageArgs struct {
	StartTime int64 `query:"start_time"`
	EndTime   int64 `query:"end_time"`
}

type CPUUsageReply struct {
	Data []Usage `json:"data"`
}

type MemoryArgs struct{}

type MemoryInfoReply struct {
	Percent float64 `json:"percent"`
	Total   float64 `json:"total"`
	Used    float64 `json:"used"`
}

type MemoryUsageArgs struct {
	StartTime int64 `query:"start_time"`
	EndTime   int64 `query:"end_time"`
}

type MemoryUsageReply struct {
	Data []Usage `json:"data"`
}

type DiskArgs struct{}

type DiskInfo struct {
	Device  string  `json:"device"`
	Percent float64 `json:"percent"`
	Total   float64 `json:"total"`
	Used    float64 `json:"used"`
}

type DiskUsage struct {
	Device string   `json:"device"`
	Data   []DiskIO `json:"data"`
}

type NetUsage struct {
	Ethernet string  `json:"ethernet"`
	Data     []NetIO `json:"data"`
}

type DiskInfoReply struct {
	Info []DiskInfo `json:"info"`
}

type DiskUsageArgs struct {
	StartTime int64 `query:"start_time"`
	EndTime   int64 `query:"end_time"`
}

type DiskUsageReply struct {
	Usage []DiskUsage `json:"usage"`
}

type NetworkUsageArgs struct {
	StartTime int64 `query:"start_time"`
	EndTime   int64 `query:"end_time"`
}

type NetworkUsageReply struct {
	Usage []NetUsage `json:"usage"`
}

type FilesSearchArgs struct {
	Path string `json:"path" validate:"required"`
}
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Mode    string `json:"mode"`
	ModTime int64  `json:"mod_time"`
	IsDir   bool   `json:"is_dir"`
}
type FilesSearchReply struct {
	Files []FileInfo `json:"files"`
}

type FileUploadArgs struct {
	SourceFilePath string `json:"source_file_path"`
	TargetFilePath string `json:"target_file_path"`
}

type FileUploadReply struct{}

type FileDownloadArgs struct {
	Filepath string `json:"filepath" validate:"required"`
}

type FileRemoteDownloadArgs struct {
	SourceFilePath string `json:"source_file_path"`
	TargetFilePath string `json:"target_file_path"`
}

type FileDownloadReply struct {
	Filepath string
}

type FileDeleteArgs struct {
	Filepath string `json:"filepath" validate:"required"`
}

type FileDeleteReply struct{}

type FileCreateArgs struct {
	Path     string `json:"path" validate:"required"`
	FileName string `json:"file_name" validate:"required"`
}

type FileCreateReply struct{}

type FolderCreateArgs struct {
	Path       string `json:"path" validate:"required"`
	FolderName string `json:"folder_name" validate:"required"`
}

type FolderCreateReply struct{}

type GetDNSSettingsArgs struct{}

type GetDNSSettingsReply struct {
	DNS []string `json:"dns"`
}

type SetDNSSettingsArgs struct {
	DNS []string `json:"dns"`
}

type SetDNSSettingsReply struct{}

type GetSystemTimeArgs struct {
}

type GetSystemTimeReply struct {
	SystemTime string `json:"system_time"`
}

type SetSystemTimeArgs struct {
	SystemTime string `json:"system_time" validate:"required"`
}

type SetSystemTimeReply struct{}

type GetSystemTimeZoneListArgs struct {
}

type GetSystemTimeZoneListReply struct {
	SystemTimeZoneList []string `json:"system_timezone_list"`
}

type GetSystemTimeZoneArgs struct{}

type GetSystemTimeZoneReply struct {
	SystemTimeZone string `json:"system_timezone"`
}

type SetSystemTimeZoneArgs struct {
	SystemTimeZone string `json:"system_timezone" validate:"required"`
}

type SetSystemTimeZoneReply struct{}

type RebootArgs struct{}

type RebootReply struct{}

type ShutdownArgs struct{}

type ShutdownReply struct{}
