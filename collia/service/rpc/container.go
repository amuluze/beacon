// Package rpc
// Container, image, and network operational RPC methods.
// Monitoring queries have been removed; data is pushed via the report mechanism.
package rpc

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	rpcSchema "common/rpc/schema"
	tunnel "common/rpc/tunnel"

	"collia/pkg/utils"
	"collia/service/model"

	"github.com/amuluze/docker"
)

// ── Container operations ──

func (s *Service) ContainerCreate(ctx context.Context, args rpcSchema.ContainerCreateArgs, reply *rpcSchema.ContainerCreateReply) error {
	var count int64
	if err := s.DB.Model(&model.Container{}).Where("name = ?", args.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("container %s already exists", args.Name)
	}
	var containerID string
	var err error
	if containerID, err = s.Manager.CreateContainer(
		ctx,
		args.Name,
		args.Image,
		args.Network,
		args.Ports,
		args.Volumes,
		args.Environments,
		nil,
		args.Labels,
	); err != nil {
		return err
	}
	updater := s.restartPolicyUpdater
	if updater == nil {
		updater = dockerRestartPolicyUpdater{}
	}
	if err := updater.UpdateRestartPolicy(ctx, containerID, normalizeRestartPolicy(args.RestartPolicy)); err != nil {
		_ = s.Manager.DeleteContainer(ctx, containerID)
		return fmt.Errorf("apply restart policy: %w", err)
	}
	reply.ContainerID = containerID
	return nil
}

func (s *Service) ContainerUpdate(ctx context.Context, args rpcSchema.ContainerUpdateArgs, reply *rpcSchema.ContainerUpdateReply) error {
	updater := s.restartPolicyUpdater
	if updater == nil {
		updater = dockerRestartPolicyUpdater{}
	}
	containerID, err := recreateContainer(ctx, s.Manager, updater, args)
	if err != nil {
		return err
	}
	reply.ContainerID = containerID
	return nil
}

func (s *Service) ContainerDelete(ctx context.Context, args rpcSchema.ContainerDeleteArgs, reply *rpcSchema.ContainerDeleteReply) error {
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Delete(&model.Container{}).Error; err != nil {
		return err
	}
	if err := s.Manager.DeleteContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerStart(ctx context.Context, args rpcSchema.ContainerStartArgs, reply *rpcSchema.ContainerStartReply) error {
	if err := s.Manager.StartContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "running").Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerStop(ctx context.Context, args rpcSchema.ContainerStopArgs, reply *rpcSchema.ContainerStopReply) error {
	if err := s.Manager.StopContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "stopped").Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerRestart(ctx context.Context, args rpcSchema.ContainerRestartArgs, reply *rpcSchema.ContainerRestartReply) error {
	if err := s.Manager.RestartContainer(ctx, args.ContainerID); err != nil {
		return err
	}
	if err := s.DB.Model(&model.Container{}).Where("container_id = ?", args.ContainerID).Update("state", "running").Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) ContainerLogs(ctx context.Context, args rpcSchema.ContainerLogsArgs, reply *rpcSchema.ContainerLogsReply) error {
	// Use streaming implementation for new pipeline; legacy stub kept for interface compliance
	return nil
}

func (s *Service) ContainerLogsStream(ctx context.Context, args rpcSchema.ContainerLogsArgs, streamSender func(frame *tunnel.Frame)) error {
	reader, err := s.Manager.ContainerLogs(ctx, args.ContainerID)
	if err != nil {
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > 8 {
			line = line[8:]
		}
		// Send each line as a stream chunk
		frame := &tunnel.Frame{
			Payload: append(line, '\n'),
		}
		streamSender(frame)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	// Signal end of stream
	streamSender(&tunnel.Frame{Eos: true})
	return nil
}

// ── Image operations ──

func (s *Service) ImagePull(ctx context.Context, args rpcSchema.ImagePullArgs, reply *rpcSchema.ImagePullReply) error {
	if err := s.Manager.PullImage(ctx, args.ImageName); err != nil {
		return err
	}
	return nil
}

func (s *Service) ImageTag(ctx context.Context, args rpcSchema.ImageTagArgs, reply *rpcSchema.ImageTagReply) error {
	if err := s.Manager.TagImage(ctx, args.OldTag, args.NewTag); err != nil {
		return err
	}
	return nil
}

func (s *Service) ImageDelete(ctx context.Context, args rpcSchema.ImageDeleteArgs, reply *rpcSchema.ImageDeleteReply) error {
	if err := s.Manager.DeleteImage(ctx, args.ImageID); err != nil {
		return err
	}
	return nil
}

func (s *Service) ImagesPrune(ctx context.Context) error {
	return s.Manager.PruneImages(ctx)
}

func (s *Service) ImageImport(ctx context.Context, args rpcSchema.ImageImportArgs, reply *rpcSchema.ImageImportReply) error {
	sourceFile := args.SourceFile
	if len(args.Data) > 0 {
		fileName := args.FileName
		if fileName == "" {
			fileName = "image.tar"
		}
		tmpFile, err := os.CreateTemp("", "collia-image-import-*-"+filepath.Base(fileName))
		if err != nil {
			return err
		}
		sourceFile = tmpFile.Name()
		if _, err := tmpFile.Write(args.Data); err != nil {
			_ = tmpFile.Close()
			_ = os.Remove(sourceFile)
			return err
		}
		if err := tmpFile.Close(); err != nil {
			_ = os.Remove(sourceFile)
			return err
		}
		defer os.Remove(sourceFile)
	}
	if err := s.Manager.ImportImage(ctx, sourceFile); err != nil {
		return err
	}
	return nil
}

func (s *Service) ImageExport(ctx context.Context, args rpcSchema.ImageExportArgs, reply *rpcSchema.ImageExportReply) error {
	targetFile := args.TargetFile
	if targetFile == "" {
		tmpFile, err := os.CreateTemp("", "collia-image-export-*.tar")
		if err != nil {
			return err
		}
		targetFile = tmpFile.Name()
		if err := tmpFile.Close(); err != nil {
			_ = os.Remove(targetFile)
			return err
		}
		defer os.Remove(targetFile)
	}
	if err := s.Manager.ExportImage(ctx, args.ImageIDs, targetFile); err != nil {
		return err
	}
	safeTarget, err := utils.SanitizePath(targetFile)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(safeTarget) //#nosec G304 -- path sanitized by utils.SanitizePath
	if err != nil && err != io.EOF {
		return err
	}
	reply.FileName = filepath.Base(targetFile)
	reply.Data = data
	return nil
}

// ── Network operations ──

func (s *Service) NetworkCreate(ctx context.Context, args rpcSchema.NetworkCreateArgs, reply *rpcSchema.NetworkCreateReply) error {
	if args.Labels == nil {
		args.Labels = make(map[string]string)
	}
	args.Labels[docker.CreatedByProbe] = "true"
	if networkID, err := s.Manager.CreateNetwork(ctx, args.Name, args.Driver, args.Subnet, args.Gateway, args.Labels); err != nil {
		return err
	} else {
		reply.NetworkID = networkID
		return nil
	}
}

func (s *Service) NetworkDelete(ctx context.Context, args rpcSchema.NetworkDeleteArgs, reply *rpcSchema.NetworkDeleteReply) error {
	if err := s.Manager.DeleteNetwork(ctx, args.NetworkID); err != nil {
		return err
	}
	return nil
}

// ── Registration ──

func registerContainerHandlers(d *Dispatcher, svc *Service) {
	// Container operations
	RegisterUnary[rpcSchema.ContainerCreateArgs, rpcSchema.ContainerCreateReply](d, "ContainerCreate", svc.ContainerCreate)
	RegisterUnary[rpcSchema.ContainerUpdateArgs, rpcSchema.ContainerUpdateReply](d, "ContainerUpdate", svc.ContainerUpdate)
	RegisterUnary[rpcSchema.ContainerDeleteArgs, rpcSchema.ContainerDeleteReply](d, "ContainerDelete", svc.ContainerDelete)
	RegisterUnary[rpcSchema.ContainerStartArgs, rpcSchema.ContainerStartReply](d, "ContainerStart", svc.ContainerStart)
	RegisterUnary[rpcSchema.ContainerStopArgs, rpcSchema.ContainerStopReply](d, "ContainerStop", svc.ContainerStop)
	RegisterUnary[rpcSchema.ContainerRestartArgs, rpcSchema.ContainerRestartReply](d, "ContainerRestart", svc.ContainerRestart)
	RegisterStream[rpcSchema.ContainerLogsArgs](d, "ContainerLogs", svc.ContainerLogsStream)

	// Image operations
	RegisterUnary[rpcSchema.ImagePullArgs, rpcSchema.ImagePullReply](d, "ImagePull", svc.ImagePull)
	RegisterUnary[rpcSchema.ImageTagArgs, rpcSchema.ImageTagReply](d, "ImageTag", svc.ImageTag)
	RegisterUnary[rpcSchema.ImageDeleteArgs, rpcSchema.ImageDeleteReply](d, "ImageDelete", svc.ImageDelete)
	d.Register("ImagesPrune", func(ctx context.Context, payload []byte, _ func(*tunnel.Frame)) ([]byte, error) {
		if err := svc.ImagesPrune(ctx); err != nil {
			return nil, err
		}
		return json.Marshal(struct{}{})
	})
	RegisterUnary[rpcSchema.ImageImportArgs, rpcSchema.ImageImportReply](d, "ImageImport", svc.ImageImport)
	RegisterUnary[rpcSchema.ImageExportArgs, rpcSchema.ImageExportReply](d, "ImageExport", svc.ImageExport)

	// Network operations
	RegisterUnary[rpcSchema.NetworkCreateArgs, rpcSchema.NetworkCreateReply](d, "NetworkCreate", svc.NetworkCreate)
	RegisterUnary[rpcSchema.NetworkDeleteArgs, rpcSchema.NetworkDeleteReply](d, "NetworkDelete", svc.NetworkDelete)
}
