// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package rpc

import (
	"common/rpc/schema"
	"context"
)

type IService interface {
	ContainerCreate(context.Context, schema.ContainerCreateArgs, *schema.ContainerCreateReply) error
	ContainerUpdate(context.Context, schema.ContainerUpdateArgs, *schema.ContainerUpdateReply) error
	ContainerDelete(context.Context, schema.ContainerDeleteArgs, *schema.ContainerDeleteReply) error
	ContainerStart(context.Context, schema.ContainerStartArgs, *schema.ContainerStartReply) error
	ContainerStop(context.Context, schema.ContainerStopArgs, *schema.ContainerStopReply) error
	ContainerRestart(context.Context, schema.ContainerRestartArgs, *schema.ContainerRestartReply) error
	ImagePull(context.Context, schema.ImagePullArgs, *schema.ImagePullReply) error
	ImageTag(context.Context, schema.ImageTagArgs, *schema.ImageTagReply) error
	ImageDelete(context.Context, schema.ImageDeleteArgs, *schema.ImageDeleteReply) error
	ImagesPrune(ctx context.Context) error
	ImageImport(context.Context, schema.ImageImportArgs, *schema.ImageImportReply) error
	ImageExport(context.Context, schema.ImageExportArgs, *schema.ImageExportReply) error
	NetworkCreate(context.Context, schema.NetworkCreateArgs, *schema.NetworkCreateReply) error
	NetworkDelete(context.Context, schema.NetworkDeleteArgs, *schema.NetworkDeleteReply) error

	FilesSearch(context.Context, schema.FilesSearchArgs, *schema.FilesSearchReply) error
	DirSize(context.Context, schema.DirSizeArgs, *schema.DirSizeReply) error
	FileCreate(context.Context, schema.FileCreateArgs, *schema.FileCreateReply) error
	FileDelete(context.Context, schema.FileDeleteArgs, *schema.FileDeleteReply) error
	FileUpload(context.Context, schema.FileUploadArgs, *schema.FileUploadReply) error
	FileDownload(context.Context, schema.FileDownloadArgs, *schema.FileDownloadReply) error
	FolderCreate(context.Context, schema.FolderCreateArgs, *schema.FolderCreateReply) error
	Reboot(context.Context, schema.RebootArgs, *schema.RebootReply) error
	Shutdown(context.Context, schema.ShutdownArgs, *schema.ShutdownReply) error
	GetDNS(context.Context, schema.GetDNSArgs, *schema.GetDNSReply) error
	SetDNS(context.Context, schema.SetDNSArgs, *schema.SetDNSReply) error
	GetSystemTime(context.Context, schema.GetSystemTimeArgs, *schema.GetSystemTimeReply) error
	SetSystemTime(context.Context, schema.SetSystemTimeArgs, *schema.SetSystemTimeReply) error
	GetSystemTimeZone(context.Context, schema.GetSystemTimeZoneArgs, *schema.GetSystemTimeZoneReply) error
	SetSystemTimeZone(context.Context, schema.SetSystemTimeZoneArgs, *schema.SetSystemTimeZoneReply) error
	GetDockerRegistryMirrors(context.Context, schema.GetDockerRegistryMirrorsArgs, *schema.GetDockerRegistryMirrorsReply) error
	SetDockerRegistryMirrors(context.Context, schema.SetDockerRegistryMirrorsArgs, *schema.SetDockerRegistryMirrorsReply) error
}
