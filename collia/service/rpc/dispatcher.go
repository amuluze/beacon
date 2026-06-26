// Package rpc
// Date: 2024/06/25
// Author: Amu
// Description: RPC dispatcher for reverse tunnel. Routes method name + JSON args
// to the appropriate Service method.
package rpc

import (
	"context"
	"encoding/json"

	rpcSchema "common/rpc/schema"
	tunnel "common/rpc/tunnel"
)

// Dispatcher routes incoming RPC calls to Service methods.
type Dispatcher struct {
	svc *Service
}

// NewDispatcher creates a new dispatcher for the given service.
func NewDispatcher(svc *Service) *Dispatcher {
	return &Dispatcher{svc: svc}
}

// Call dispatches the method call with JSON-encoded payload and returns JSON-encoded response.
// For streaming methods (ContainerLogs), streamSender sends multiple data frames.
func (d *Dispatcher) Call(ctx context.Context, method string, payload []byte, streamSender func(*tunnel.Frame)) ([]byte, error) {
	switch method {
	// ── Container operations ──
	case "ContainerCreate":
		var args rpcSchema.ContainerCreateArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ContainerCreateReply
		if err := d.svc.ContainerCreate(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ContainerUpdate":
		var args rpcSchema.ContainerUpdateArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ContainerUpdateReply
		if err := d.svc.ContainerUpdate(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ContainerDelete":
		var args rpcSchema.ContainerDeleteArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ContainerDeleteReply
		if err := d.svc.ContainerDelete(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ContainerStart":
		var args rpcSchema.ContainerStartArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ContainerStartReply
		if err := d.svc.ContainerStart(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ContainerStop":
		var args rpcSchema.ContainerStopArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ContainerStopReply
		if err := d.svc.ContainerStop(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ContainerRestart":
		var args rpcSchema.ContainerRestartArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ContainerRestartReply
		if err := d.svc.ContainerRestart(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ContainerLogs":
		var args rpcSchema.ContainerLogsArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		err := d.svc.ContainerLogsStream(ctx, args, streamSender)
		return nil, err

	// ── Image operations ──
	case "ImagePull":
		var args rpcSchema.ImagePullArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ImagePullReply
		if err := d.svc.ImagePull(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ImageTag":
		var args rpcSchema.ImageTagArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ImageTagReply
		if err := d.svc.ImageTag(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ImageDelete":
		var args rpcSchema.ImageDeleteArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ImageDeleteReply
		if err := d.svc.ImageDelete(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ImagesPrune":
		if err := d.svc.ImagesPrune(ctx); err != nil {
			return nil, err
		}
		return json.Marshal(struct{}{})
	case "ImageImport":
		var args rpcSchema.ImageImportArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ImageImportReply
		if err := d.svc.ImageImport(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "ImageExport":
		var args rpcSchema.ImageExportArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ImageExportReply
		if err := d.svc.ImageExport(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)

	// ── Network operations ──
	case "NetworkCreate":
		var args rpcSchema.NetworkCreateArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.NetworkCreateReply
		if err := d.svc.NetworkCreate(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "NetworkDelete":
		var args rpcSchema.NetworkDeleteArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.NetworkDeleteReply
		if err := d.svc.NetworkDelete(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)

	// ── File operations ──
	case "FilesSearch":
		var args rpcSchema.FilesSearchArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.FilesSearchReply
		if err := d.svc.FilesSearch(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "DirSize":
		var args rpcSchema.DirSizeArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.DirSizeReply
		if err := d.svc.DirSize(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "FileCreate":
		var args rpcSchema.FileCreateArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.FileCreateReply
		if err := d.svc.FileCreate(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "FileDelete":
		var args rpcSchema.FileDeleteArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.FileDeleteReply
		if err := d.svc.FileDelete(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "FileUpload":
		var args rpcSchema.FileUploadArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.FileUploadReply
		if err := d.svc.FileUpload(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "FileDownload":
		var args rpcSchema.FileDownloadArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.FileDownloadReply
		if err := d.svc.FileDownload(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "FolderCreate":
		var args rpcSchema.FolderCreateArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.FolderCreateReply
		if err := d.svc.FolderCreate(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)

	// ── System operations ──
	case "Reboot":
		var args rpcSchema.RebootArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.RebootReply
		if err := d.svc.Reboot(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "Shutdown":
		var args rpcSchema.ShutdownArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ShutdownReply
		if err := d.svc.Shutdown(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "GetDNS":
		var args rpcSchema.GetDNSArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.GetDNSReply
		if err := d.svc.GetDNS(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "SetDNS":
		var args rpcSchema.SetDNSArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.SetDNSReply
		if err := d.svc.SetDNS(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "GetSystemTime":
		var args rpcSchema.GetSystemTimeArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.GetSystemTimeReply
		if err := d.svc.GetSystemTime(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "SetSystemTime":
		var args rpcSchema.SetSystemTimeArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.SetSystemTimeReply
		if err := d.svc.SetSystemTime(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "GetSystemTimeZone":
		var args rpcSchema.GetSystemTimeZoneArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.GetSystemTimeZoneReply
		if err := d.svc.GetSystemTimeZone(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "SetSystemTimeZone":
		var args rpcSchema.SetSystemTimeZoneArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.SetSystemTimeZoneReply
		if err := d.svc.SetSystemTimeZone(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "GetSystemTimeZoneList":
		var args rpcSchema.GetSystemTimeZoneListArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.GetSystemTimeZoneListReply
		if err := d.svc.GetSystemTimeZoneList(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "GetDockerRegistryMirrors":
		var args rpcSchema.GetDockerRegistryMirrorsArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.GetDockerRegistryMirrorsReply
		if err := d.svc.GetDockerRegistryMirrors(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "SetDockerRegistryMirrors":
		var args rpcSchema.SetDockerRegistryMirrorsArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.SetDockerRegistryMirrorsReply
		if err := d.svc.SetDockerRegistryMirrors(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)

	// ── Terminal operations ──
	case "TerminalSession":
		var args rpcSchema.TerminalSessionArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		err := d.svc.TerminalSessionStream(ctx, args, streamSender)
		return nil, err
	case "ResizeTerminal":
		var args rpcSchema.ResizeTerminalArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.ResizeTerminalReply
		if err := d.svc.ResizeTerminal(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "TerminalInput":
		var args rpcSchema.TerminalInputArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.TerminalInputReply
		if err := d.svc.TerminalInput(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	case "TerminalClose":
		var args rpcSchema.TerminalCloseArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, err
		}
		var reply rpcSchema.TerminalCloseReply
		if err := d.svc.TerminalClose(ctx, args, &reply); err != nil {
			return nil, err
		}
		return json.Marshal(reply)

	default:
		return nil, &UnknownMethodError{Method: method}
	}
}

// UnknownMethodError is returned when a method is not recognized.
type UnknownMethodError struct {
	Method string
}

func (e *UnknownMethodError) Error() string {
	return "unknown rpc method: " + e.Method
}
