// Package testutil
// Date: 2026/6/26
// Author: Amu
// Description: shared fake implementations for service layer unit tests
package testutil

import (
	"context"
	"fmt"

	"amprobe/pkg/auth"
	alarmRepo "amprobe/service/alarm/repository"
	auditRepo "amprobe/service/audit/repository"
	authRepo "amprobe/service/auth/repository"
	accountRepo "amprobe/service/account/repository"
	mailRepo "amprobe/service/mail/repository"
	"amprobe/service/model"
	"amprobe/service/schema"

	rpcSchema "common/rpc/schema"

	"gorm.io/gorm"
)

// ── FakeAuther ──

type FakeAuther struct {
	GenerateTokenFn func(userID, username string) (auth.TokenInfo, error)
	ParseTokenFn    func(token, tokenType string) (string, string, error)
	DestroyTokenFn  func(token string) error
	ReleaseFn       func() error
	RecordAuditFn   func(username, operate string)
}

func NewFakeAuther() *FakeAuther {
	return &FakeAuther{
		GenerateTokenFn: func(userID, username string) (auth.TokenInfo, error) {
			return &fakeTokenInfo{accessToken: "at-" + userID, refreshToken: "rt-" + userID}, nil
		},
		ParseTokenFn: func(token, tokenType string) (string, string, error) {
			return "user-1", "admin", nil
		},
		DestroyTokenFn: func(token string) error { return nil },
		ReleaseFn:      func() error { return nil },
		RecordAuditFn:  func(string, string) {},
	}
}

func (f *FakeAuther) GenerateToken(userID, username string) (auth.TokenInfo, error) {
	return f.GenerateTokenFn(userID, username)
}
func (f *FakeAuther) ParseToken(token, tokenType string) (string, string, error) {
	return f.ParseTokenFn(token, tokenType)
}
func (f *FakeAuther) DestroyToken(token string) error { return f.DestroyTokenFn(token) }
func (f *FakeAuther) Release() error                  { return f.ReleaseFn() }
func (f *FakeAuther) RecordAudit(username, operate string) { f.RecordAuditFn(username, operate) }

type fakeTokenInfo struct {
	accessToken  string
	refreshToken string
}

func (t *fakeTokenInfo) GetAccessToken() string  { return t.accessToken }
func (t *fakeTokenInfo) GetRefreshToken() string { return t.refreshToken }

// ── FakeAuthRepo ──

type FakeAuthRepo struct {
	LoginFn    func(ctx context.Context, args schema.LoginArgs) (model.User, error)
	PassUpdateFn func(ctx context.Context, args schema.PasswordUpdateArgs) error
	UserInfoFn  func(ctx context.Context, userID string) (model.User, error)
}

func NewFakeAuthRepo() *FakeAuthRepo {
	return &FakeAuthRepo{
		LoginFn: func(ctx context.Context, args schema.LoginArgs) (model.User, error) {
			return model.User{Username: args.Username}, nil
		},
		PassUpdateFn: func(ctx context.Context, args schema.PasswordUpdateArgs) error { return nil },
		UserInfoFn: func(ctx context.Context, userID string) (model.User, error) {
			return model.User{Username: "admin"}, nil
		},
	}
}

func (f *FakeAuthRepo) Login(ctx context.Context, args schema.LoginArgs) (model.User, error) {
	return f.LoginFn(ctx, args)
}
func (f *FakeAuthRepo) PassUpdate(ctx context.Context, args schema.PasswordUpdateArgs) error {
	return f.PassUpdateFn(ctx, args)
}
func (f *FakeAuthRepo) UserInfo(ctx context.Context, userID string) (model.User, error) {
	return f.UserInfoFn(ctx, userID)
}

// Verify interface compliance
var _ authRepo.IAuthRepository = (*FakeAuthRepo)(nil)

// ── FakeAccountRepo ──

type FakeAccountRepo struct {
	UserQueryFn    func(ctx context.Context, args schema.UserQueryArgs) (model.Users, error)
	UserCreateFn   func(ctx context.Context, args schema.UserCreateArgs) (model.User, error)
	UserUpdateFn   func(ctx context.Context, args schema.UserUpdateArgs) (model.User, error)
	UserDeleteFn   func(ctx context.Context, args schema.UserDeleteArgs) error
	UserCountFn    func(ctx context.Context) (int64, error)
	RoleQueryFn    func(ctx context.Context, args schema.RoleQueryArgs) (model.Roles, error)
	RoleCreateFn   func(ctx context.Context, args schema.RoleCreateArgs) (model.Role, error)
	RoleUpdateFn   func(ctx context.Context, args schema.RoleUpdateArgs) (model.Role, error)
	RoleDeleteFn   func(ctx context.Context, args schema.RoleDeleteArgs) error
	RoleCountFn    func(ctx context.Context) (int64, error)
	ResourceQueryFn func(ctx context.Context, args schema.ResourceQueryArgs) (model.Resources, error)
	ResourceCountFn func(ctx context.Context) (int64, error)
}

func NewFakeAccountRepo() *FakeAccountRepo {
	return &FakeAccountRepo{
		UserQueryFn:    func(ctx context.Context, args schema.UserQueryArgs) (model.Users, error) { return nil, nil },
		UserCreateFn:   func(ctx context.Context, args schema.UserCreateArgs) (model.User, error) { return model.User{}, nil },
		UserUpdateFn:   func(ctx context.Context, args schema.UserUpdateArgs) (model.User, error) { return model.User{}, nil },
		UserDeleteFn:   func(ctx context.Context, args schema.UserDeleteArgs) error { return nil },
		UserCountFn:    func(ctx context.Context) (int64, error) { return 0, nil },
		RoleQueryFn:    func(ctx context.Context, args schema.RoleQueryArgs) (model.Roles, error) { return nil, nil },
		RoleCreateFn:   func(ctx context.Context, args schema.RoleCreateArgs) (model.Role, error) { return model.Role{}, nil },
		RoleUpdateFn:   func(ctx context.Context, args schema.RoleUpdateArgs) (model.Role, error) { return model.Role{}, nil },
		RoleDeleteFn:   func(ctx context.Context, args schema.RoleDeleteArgs) error { return nil },
		RoleCountFn:    func(ctx context.Context) (int64, error) { return 0, nil },
		ResourceQueryFn: func(ctx context.Context, args schema.ResourceQueryArgs) (model.Resources, error) { return nil, nil },
		ResourceCountFn: func(ctx context.Context) (int64, error) { return 0, nil },
	}
}

func (f *FakeAccountRepo) UserQuery(ctx context.Context, args schema.UserQueryArgs) (model.Users, error) {
	return f.UserQueryFn(ctx, args)
}
func (f *FakeAccountRepo) UserCreate(ctx context.Context, args schema.UserCreateArgs) (model.User, error) {
	return f.UserCreateFn(ctx, args)
}
func (f *FakeAccountRepo) UserUpdate(ctx context.Context, args schema.UserUpdateArgs) (model.User, error) {
	return f.UserUpdateFn(ctx, args)
}
func (f *FakeAccountRepo) UserDelete(ctx context.Context, args schema.UserDeleteArgs) error {
	return f.UserDeleteFn(ctx, args)
}
func (f *FakeAccountRepo) UserCount(ctx context.Context) (int64, error) { return f.UserCountFn(ctx) }
func (f *FakeAccountRepo) RoleQuery(ctx context.Context, args schema.RoleQueryArgs) (model.Roles, error) {
	return f.RoleQueryFn(ctx, args)
}
func (f *FakeAccountRepo) RoleCreate(ctx context.Context, args schema.RoleCreateArgs) (model.Role, error) {
	return f.RoleCreateFn(ctx, args)
}
func (f *FakeAccountRepo) RoleUpdate(ctx context.Context, args schema.RoleUpdateArgs) (model.Role, error) {
	return f.RoleUpdateFn(ctx, args)
}
func (f *FakeAccountRepo) RoleDelete(ctx context.Context, args schema.RoleDeleteArgs) error {
	return f.RoleDeleteFn(ctx, args)
}
func (f *FakeAccountRepo) RoleCount(ctx context.Context) (int64, error) { return f.RoleCountFn(ctx) }
func (f *FakeAccountRepo) ResourceQuery(ctx context.Context, args schema.ResourceQueryArgs) (model.Resources, error) {
	return f.ResourceQueryFn(ctx, args)
}
func (f *FakeAccountRepo) ResourceCount(ctx context.Context) (int64, error) { return f.ResourceCountFn(ctx) }

var _ accountRepo.IAccountRepository = (*FakeAccountRepo)(nil)

// ── FakeAlarmRepo ──

type FakeAlarmRepo struct {
	AlarmQueryFn  func(ctx context.Context) ([]model.AlarmThreshold, error)
	AlarmUpdateFn func(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error
}

func NewFakeAlarmRepo() *FakeAlarmRepo {
	return &FakeAlarmRepo{
		AlarmQueryFn:  func(ctx context.Context) ([]model.AlarmThreshold, error) { return nil, nil },
		AlarmUpdateFn: func(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error { return nil },
	}
}

func (f *FakeAlarmRepo) AlarmQuery(ctx context.Context) ([]model.AlarmThreshold, error) {
	return f.AlarmQueryFn(ctx)
}
func (f *FakeAlarmRepo) AlarmUpdate(ctx context.Context, args schema.AlarmThresholdUpdateArgs) error {
	return f.AlarmUpdateFn(ctx, args)
}

var _ alarmRepo.IAlarmRepository = (*FakeAlarmRepo)(nil)

// ── FakeMailRepo ──

type FakeMailRepo struct {
	MailQueryFn  func(ctx context.Context) (model.Mail, error)
	MailCreateFn func(ctx context.Context, args schema.MailCreateArgs) error
	MailUpdateFn func(ctx context.Context, args schema.MailUpdateArgs) error
	MailDeleteFn func(ctx context.Context, args schema.MailDeleteArgs) error
	MailTestFn   func(ctx context.Context, args schema.MailTestArgs) error
}

func NewFakeMailRepo() *FakeMailRepo {
	return &FakeMailRepo{
		MailQueryFn:  func(ctx context.Context) (model.Mail, error) { return model.Mail{}, nil },
		MailCreateFn: func(ctx context.Context, args schema.MailCreateArgs) error { return nil },
		MailUpdateFn: func(ctx context.Context, args schema.MailUpdateArgs) error { return nil },
		MailDeleteFn: func(ctx context.Context, args schema.MailDeleteArgs) error { return nil },
		MailTestFn:   func(ctx context.Context, args schema.MailTestArgs) error { return nil },
	}
}

func (f *FakeMailRepo) MailQuery(ctx context.Context) (model.Mail, error)   { return f.MailQueryFn(ctx) }
func (f *FakeMailRepo) MailCreate(ctx context.Context, args schema.MailCreateArgs) error {
	return f.MailCreateFn(ctx, args)
}
func (f *FakeMailRepo) MailUpdate(ctx context.Context, args schema.MailUpdateArgs) error {
	return f.MailUpdateFn(ctx, args)
}
func (f *FakeMailRepo) MailDelete(ctx context.Context, args schema.MailDeleteArgs) error {
	return f.MailDeleteFn(ctx, args)
}
func (f *FakeMailRepo) MailTest(ctx context.Context, args schema.MailTestArgs) error {
	return f.MailTestFn(ctx, args)
}

var _ mailRepo.IMailRepository = (*FakeMailRepo)(nil)

// ── FakeAuditRepo ──

type FakeAuditRepo struct {
	AuditQueryFn func(ctx context.Context, args schema.AuditQueryArgs) (model.Audits, error)
	AuditCountFn func(ctx context.Context) (int, error)
}

func NewFakeAuditRepo() *FakeAuditRepo {
	return &FakeAuditRepo{
		AuditQueryFn: func(ctx context.Context, args schema.AuditQueryArgs) (model.Audits, error) { return nil, nil },
		AuditCountFn: func(ctx context.Context) (int, error) { return 0, nil },
	}
}

func (f *FakeAuditRepo) AuditQuery(ctx context.Context, args schema.AuditQueryArgs) (model.Audits, error) {
	return f.AuditQueryFn(ctx, args)
}
func (f *FakeAuditRepo) AuditCount(ctx context.Context) (int, error) { return f.AuditCountFn(ctx) }

var _ auditRepo.IAuditRepo = (*FakeAuditRepo)(nil)

// ── ErrTest ──

var ErrTest = fmt.Errorf("test error")

// ── Stubs for rpc schema types used in host/container tests ──
// These provide zero-value returns to satisfy interfaces.

// FakeHostRepo implements host/repository.IHostRepo
type FakeHostRepo struct {
	HostInfoFn    func(context.Context, rpcSchema.HostInfoArgs) (rpcSchema.HostInfoReply, error)
	CPUInfoFn     func(context.Context, rpcSchema.CPUInfoArgs) (rpcSchema.CPUInfoReply, error)
	CPUUsageFn    func(context.Context, rpcSchema.CPUUsageArgs) (rpcSchema.CPUUsageReply, error)
	MemInfoFn     func(context.Context, rpcSchema.MemoryInfoArgs) (rpcSchema.MemoryInfoReply, error)
	MemUsageFn    func(context.Context, rpcSchema.MemoryUsageArgs) (rpcSchema.MemoryUsageReply, error)
	DiskInfoFn    func(context.Context, rpcSchema.DiskInfoArgs) (rpcSchema.DiskInfoReply, error)
	DiskUsageFn   func(context.Context, rpcSchema.DiskUsageArgs) (rpcSchema.DiskUsageReply, error)
	NetUsageFn    func(context.Context, rpcSchema.NetUsageArgs) (rpcSchema.NetUsageReply, error)
	FilesSearchFn func(context.Context, rpcSchema.FilesSearchArgs) (rpcSchema.FilesSearchReply, error)
	FileUploadFn  func(context.Context, rpcSchema.FileUploadArgs) error
	FileDownloadFn func(context.Context, rpcSchema.FileDownloadArgs) (rpcSchema.FileDownloadReply, error)
	FileDeleteFn  func(context.Context, rpcSchema.FileDeleteArgs) error
	FileCreateFn  func(context.Context, rpcSchema.FileCreateArgs) error
	FolderCreateFn func(context.Context, rpcSchema.FolderCreateArgs) error
	GetDNSFn      func(context.Context, rpcSchema.GetDNSArgs) (rpcSchema.GetDNSReply, error)
	SetDNSFn      func(context.Context, rpcSchema.SetDNSArgs) error
	GetTimeFn     func(context.Context, rpcSchema.GetSystemTimeArgs) (rpcSchema.GetSystemTimeReply, error)
	SetTimeFn     func(context.Context, rpcSchema.SetSystemTimeArgs) error
	GetTZListFn   func(context.Context, rpcSchema.GetSystemTimeZoneListArgs) (rpcSchema.GetSystemTimeZoneListReply, error)
	GetTZFn       func(context.Context, rpcSchema.GetSystemTimeZoneArgs) (rpcSchema.GetSystemTimeZoneReply, error)
	SetTZFn       func(context.Context, rpcSchema.SetSystemTimeZoneArgs) error
	RebootFn      func(context.Context, rpcSchema.RebootArgs) error
	ShutdownFn    func(context.Context, rpcSchema.ShutdownArgs) error
}

func (f *FakeHostRepo) HostInfo(ctx context.Context, args rpcSchema.HostInfoArgs) (rpcSchema.HostInfoReply, error) {
	return f.HostInfoFn(ctx, args)
}
func (f *FakeHostRepo) CPUInfo(ctx context.Context, args rpcSchema.CPUInfoArgs) (rpcSchema.CPUInfoReply, error) {
	return f.CPUInfoFn(ctx, args)
}
func (f *FakeHostRepo) CPUUsage(ctx context.Context, args rpcSchema.CPUUsageArgs) (rpcSchema.CPUUsageReply, error) {
	return f.CPUUsageFn(ctx, args)
}
func (f *FakeHostRepo) MemInfo(ctx context.Context, args rpcSchema.MemoryInfoArgs) (rpcSchema.MemoryInfoReply, error) {
	return f.MemInfoFn(ctx, args)
}
func (f *FakeHostRepo) MemUsage(ctx context.Context, args rpcSchema.MemoryUsageArgs) (rpcSchema.MemoryUsageReply, error) {
	return f.MemUsageFn(ctx, args)
}
func (f *FakeHostRepo) DiskInfo(ctx context.Context, args rpcSchema.DiskInfoArgs) (rpcSchema.DiskInfoReply, error) {
	return f.DiskInfoFn(ctx, args)
}
func (f *FakeHostRepo) DiskUsage(ctx context.Context, args rpcSchema.DiskUsageArgs) (rpcSchema.DiskUsageReply, error) {
	return f.DiskUsageFn(ctx, args)
}
func (f *FakeHostRepo) NetUsage(ctx context.Context, args rpcSchema.NetUsageArgs) (rpcSchema.NetUsageReply, error) {
	return f.NetUsageFn(ctx, args)
}
func (f *FakeHostRepo) FilesSearch(ctx context.Context, args rpcSchema.FilesSearchArgs) (rpcSchema.FilesSearchReply, error) {
	return f.FilesSearchFn(ctx, args)
}
func (f *FakeHostRepo) FileUpload(ctx context.Context, args rpcSchema.FileUploadArgs) error {
	return f.FileUploadFn(ctx, args)
}
func (f *FakeHostRepo) FileDownload(ctx context.Context, args rpcSchema.FileDownloadArgs) (rpcSchema.FileDownloadReply, error) {
	return f.FileDownloadFn(ctx, args)
}
func (f *FakeHostRepo) FileDelete(ctx context.Context, args rpcSchema.FileDeleteArgs) error {
	return f.FileDeleteFn(ctx, args)
}
func (f *FakeHostRepo) FileCreate(ctx context.Context, args rpcSchema.FileCreateArgs) error {
	return f.FileCreateFn(ctx, args)
}
func (f *FakeHostRepo) FolderCreate(ctx context.Context, args rpcSchema.FolderCreateArgs) error {
	return f.FolderCreateFn(ctx, args)
}
func (f *FakeHostRepo) GetDNSSettings(ctx context.Context, args rpcSchema.GetDNSArgs) (rpcSchema.GetDNSReply, error) {
	return f.GetDNSFn(ctx, args)
}
func (f *FakeHostRepo) SetDNSSettings(ctx context.Context, args rpcSchema.SetDNSArgs) error {
	return f.SetDNSFn(ctx, args)
}
func (f *FakeHostRepo) GetSystemTime(ctx context.Context, args rpcSchema.GetSystemTimeArgs) (rpcSchema.GetSystemTimeReply, error) {
	return f.GetTimeFn(ctx, args)
}
func (f *FakeHostRepo) SetSystemTime(ctx context.Context, args rpcSchema.SetSystemTimeArgs) error {
	return f.SetTimeFn(ctx, args)
}
func (f *FakeHostRepo) GetSystemTimeZoneList(ctx context.Context, args rpcSchema.GetSystemTimeZoneListArgs) (rpcSchema.GetSystemTimeZoneListReply, error) {
	return f.GetTZListFn(ctx, args)
}
func (f *FakeHostRepo) GetSystemTimeZone(ctx context.Context, args rpcSchema.GetSystemTimeZoneArgs) (rpcSchema.GetSystemTimeZoneReply, error) {
	return f.GetTZFn(ctx, args)
}
func (f *FakeHostRepo) SetSystemTimeZone(ctx context.Context, args rpcSchema.SetSystemTimeZoneArgs) error {
	return f.SetTZFn(ctx, args)
}
func (f *FakeHostRepo) Reboot(ctx context.Context, args rpcSchema.RebootArgs) error {
	return f.RebootFn(ctx, args)
}
func (f *FakeHostRepo) Shutdown(ctx context.Context, args rpcSchema.ShutdownArgs) error {
	return f.ShutdownFn(ctx, args)
}

// FakeContainerRepo implements container/repository.IContainerRepo
type FakeContainerRepo struct {
	VersionFn                func(ctx context.Context, args rpcSchema.DockerArgs) (rpcSchema.DockerReply, error)
	ContainerListFn          func(ctx context.Context, args rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error)
	UsageFn                  func(ctx context.Context, args rpcSchema.ContainerUsageArgs) (rpcSchema.ContainerUsageReply, error)
	ContainersByImageFn      func(ctx context.Context, image string) (int, error)
	ContainerCountFn         func(ctx context.Context, args rpcSchema.ContainerCountArgs) (rpcSchema.ContainerCountReply, error)
	ContainerCreateFn        func(ctx context.Context, args rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error)
	ContainerUpdateFn        func(ctx context.Context, args rpcSchema.ContainerUpdateArgs) (rpcSchema.ContainerUpdateReply, error)
	ContainerDeleteFn        func(ctx context.Context, args rpcSchema.ContainerDeleteArgs) error
	ContainerStartFn         func(ctx context.Context, args rpcSchema.ContainerStartArgs) error
	ContainerStopFn          func(ctx context.Context, args rpcSchema.ContainerStopArgs) error
	ContainerRestartFn       func(ctx context.Context, args rpcSchema.ContainerRestartArgs) error
	ContainerLogsFn          func(ctx context.Context, args rpcSchema.ContainerLogsArgs) (rpcSchema.ContainerLogsReply, error)
	ImageListFn              func(ctx context.Context, args rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error)
	ImageCountFn             func(ctx context.Context, args rpcSchema.ImageCountArgs) (rpcSchema.ImageCountReply, error)
	ImagePullFn              func(ctx context.Context, args rpcSchema.ImagePullArgs) error
	ImageTagFn               func(ctx context.Context, args rpcSchema.ImageTagArgs) error
	ImageImportFn            func(ctx context.Context, args rpcSchema.ImageImportArgs) error
	ImageExportFn            func(ctx context.Context, args rpcSchema.ImageExportArgs) (rpcSchema.ImageExportReply, error)
	ImageDeleteFn            func(ctx context.Context, args rpcSchema.ImageDeleteArgs) error
	ImagesPruneFn            func(ctx context.Context) error
	NetworkListFn            func(ctx context.Context, args rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error)
	NetworkCountFn           func(ctx context.Context, args rpcSchema.NetworkCountArgs) (rpcSchema.NetworkCountReply, error)
	NetworkCreateFn          func(ctx context.Context, args rpcSchema.NetworkCreateArgs) (rpcSchema.NetworkCreateReply, error)
	NetworkDeleteFn          func(ctx context.Context, args rpcSchema.NetworkDeleteArgs) error
	GetDockerRegistryMirrorsFn func(ctx context.Context, args rpcSchema.GetDockerRegistryMirrorsArgs) (rpcSchema.GetDockerRegistryMirrorsReply, error)
	SetDockerRegistryMirrorsFn func(ctx context.Context, args rpcSchema.SetDockerRegistryMirrorsArgs) error
}

func (f *FakeContainerRepo) Version(ctx context.Context, args rpcSchema.DockerArgs) (rpcSchema.DockerReply, error) {
	return f.VersionFn(ctx, args)
}
func (f *FakeContainerRepo) ContainerList(ctx context.Context, args rpcSchema.ContainerQueryArgs) (rpcSchema.ContainerQueryReply, error) {
	return f.ContainerListFn(ctx, args)
}
func (f *FakeContainerRepo) Usage(ctx context.Context, args rpcSchema.ContainerUsageArgs) (rpcSchema.ContainerUsageReply, error) {
	return f.UsageFn(ctx, args)
}
func (f *FakeContainerRepo) ContainersByImage(ctx context.Context, image string) (int, error) {
	return f.ContainersByImageFn(ctx, image)
}
func (f *FakeContainerRepo) ContainerCount(ctx context.Context, args rpcSchema.ContainerCountArgs) (rpcSchema.ContainerCountReply, error) {
	return f.ContainerCountFn(ctx, args)
}
func (f *FakeContainerRepo) ContainerCreate(ctx context.Context, args rpcSchema.ContainerCreateArgs) (rpcSchema.ContainerCreateReply, error) {
	return f.ContainerCreateFn(ctx, args)
}
func (f *FakeContainerRepo) ContainerUpdate(ctx context.Context, args rpcSchema.ContainerUpdateArgs) (rpcSchema.ContainerUpdateReply, error) {
	return f.ContainerUpdateFn(ctx, args)
}
func (f *FakeContainerRepo) ContainerDelete(ctx context.Context, args rpcSchema.ContainerDeleteArgs) error {
	return f.ContainerDeleteFn(ctx, args)
}
func (f *FakeContainerRepo) ContainerStart(ctx context.Context, args rpcSchema.ContainerStartArgs) error {
	return f.ContainerStartFn(ctx, args)
}
func (f *FakeContainerRepo) ContainerStop(ctx context.Context, args rpcSchema.ContainerStopArgs) error {
	return f.ContainerStopFn(ctx, args)
}
func (f *FakeContainerRepo) ContainerRestart(ctx context.Context, args rpcSchema.ContainerRestartArgs) error {
	return f.ContainerRestartFn(ctx, args)
}
func (f *FakeContainerRepo) ContainerLogs(ctx context.Context, args rpcSchema.ContainerLogsArgs) (rpcSchema.ContainerLogsReply, error) {
	return f.ContainerLogsFn(ctx, args)
}
func (f *FakeContainerRepo) ImageList(ctx context.Context, args rpcSchema.ImageQueryArgs) (rpcSchema.ImageQueryReply, error) {
	return f.ImageListFn(ctx, args)
}
func (f *FakeContainerRepo) ImageCount(ctx context.Context, args rpcSchema.ImageCountArgs) (rpcSchema.ImageCountReply, error) {
	return f.ImageCountFn(ctx, args)
}
func (f *FakeContainerRepo) ImagePull(ctx context.Context, args rpcSchema.ImagePullArgs) error {
	return f.ImagePullFn(ctx, args)
}
func (f *FakeContainerRepo) ImageTag(ctx context.Context, args rpcSchema.ImageTagArgs) error {
	return f.ImageTagFn(ctx, args)
}
func (f *FakeContainerRepo) ImageImport(ctx context.Context, args rpcSchema.ImageImportArgs) error {
	return f.ImageImportFn(ctx, args)
}
func (f *FakeContainerRepo) ImageExport(ctx context.Context, args rpcSchema.ImageExportArgs) (rpcSchema.ImageExportReply, error) {
	return f.ImageExportFn(ctx, args)
}
func (f *FakeContainerRepo) ImageDelete(ctx context.Context, args rpcSchema.ImageDeleteArgs) error {
	return f.ImageDeleteFn(ctx, args)
}
func (f *FakeContainerRepo) ImagesPrune(ctx context.Context) error { return f.ImagesPruneFn(ctx) }
func (f *FakeContainerRepo) NetworkList(ctx context.Context, args rpcSchema.NetworkQueryArgs) (rpcSchema.NetworkQueryReply, error) {
	return f.NetworkListFn(ctx, args)
}
func (f *FakeContainerRepo) NetworkCount(ctx context.Context, args rpcSchema.NetworkCountArgs) (rpcSchema.NetworkCountReply, error) {
	return f.NetworkCountFn(ctx, args)
}
func (f *FakeContainerRepo) NetworkCreate(ctx context.Context, args rpcSchema.NetworkCreateArgs) (rpcSchema.NetworkCreateReply, error) {
	return f.NetworkCreateFn(ctx, args)
}
func (f *FakeContainerRepo) NetworkDelete(ctx context.Context, args rpcSchema.NetworkDeleteArgs) error {
	return f.NetworkDeleteFn(ctx, args)
}
func (f *FakeContainerRepo) GetDockerRegistryMirrors(ctx context.Context, args rpcSchema.GetDockerRegistryMirrorsArgs) (rpcSchema.GetDockerRegistryMirrorsReply, error) {
	return f.GetDockerRegistryMirrorsFn(ctx, args)
}
func (f *FakeContainerRepo) SetDockerRegistryMirrors(ctx context.Context, args rpcSchema.SetDockerRegistryMirrorsArgs) error {
	return f.SetDockerRegistryMirrorsFn(ctx, args)
}

// suppress unused import warning
var _ = gorm.ErrRecordNotFound
