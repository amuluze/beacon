// Package rpc
// Date: 2024/06/25 00:15:52
// Author: Amu
// Description:
package rpc

import (
	"collia/pkg/utils"
	rpcSchema "common/rpc/schema"
	"context"
	"os"
	"path/filepath"
)

func (s *Service) FilesSearch(ctx context.Context, args rpcSchema.FilesSearchArgs, reply *rpcSchema.FilesSearchReply) error {
	files, err := os.ReadDir(args.Path)
	if err != nil {
		return err
	}
	data := make([]rpcSchema.FileInfo, 0)
	for _, file := range files {
		info, _ := file.Info()
		data = append(data, rpcSchema.FileInfo{
			Name:    file.Name(),
			IsDir:   file.IsDir(),
			Size:    info.Size(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime().Unix(),
		})
	}
	reply.Files = data
	return nil
}

func (s *Service) DirSize(ctx context.Context, args rpcSchema.DirSizeArgs, reply *rpcSchema.DirSizeReply) error {
	var size int64
	if err := filepath.Walk(args.Path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	}); err != nil {
		return err
	}
	reply.Size = size
	return nil
}

func (s *Service) FileCreate(ctx context.Context, args rpcSchema.FileCreateArgs, reply *rpcSchema.FileCreateReply) error {
	filePath := filepath.Join(args.Path, args.FileName)
	filePath, err := utils.SanitizePath(filePath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600) //#nosec G304 -- path sanitized by utils.SanitizePath
		if err != nil {
			return err
		}
		return f.Close()
	}
	return nil
}

func (s *Service) FolderCreate(ctx context.Context, args rpcSchema.FolderCreateArgs, reply *rpcSchema.FolderCreateReply) error {
	folderPath := filepath.Join(args.Path, args.FolderName)
	folderPath, err := utils.SanitizePath(folderPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return os.Mkdir(folderPath, 0750)
	}
	return nil
}

func (s *Service) FileDelete(ctx context.Context, args rpcSchema.FileDeleteArgs, reply *rpcSchema.FileDeleteReply) error {
	if info, err := os.Stat(args.Filepath); err != nil {
		return err
	} else if info.IsDir() {
		return os.RemoveAll(args.Filepath)
	} else {
		return os.Remove(args.Filepath)
	}
}

func (s *Service) FileUpload(ctx context.Context, args rpcSchema.FileUploadArgs, reply *rpcSchema.FileUploadReply) error {
	if len(args.Data) > 0 {
		if err := os.MkdirAll(filepath.Dir(args.TargetFilePath), 0750); err != nil {
			return err
		}
		return os.WriteFile(args.TargetFilePath, args.Data, 0600)
	}
	return os.Rename(args.SourceFilePath, args.TargetFilePath)
}

func (s *Service) FileDownload(ctx context.Context, args rpcSchema.FileDownloadArgs, reply *rpcSchema.FileDownloadReply) error {
	if args.TargetFilePath == "" {
		data, err := os.ReadFile(args.SourceFilePath)
		if err != nil {
			return err
		}
		reply.Filepath = args.SourceFilePath
		reply.FileName = filepath.Base(args.SourceFilePath)
		reply.Data = data
		return nil
	}
	_, err := utils.CopyFile(args.SourceFilePath, args.TargetFilePath)
	reply.Filepath = args.TargetFilePath
	reply.FileName = filepath.Base(args.TargetFilePath)
	return err
}

// ── Registration ──

func registerFileHandlers(d *Dispatcher, svc *Service) {
	RegisterUnary[rpcSchema.FilesSearchArgs, rpcSchema.FilesSearchReply](d, "FilesSearch", svc.FilesSearch)
	RegisterUnary[rpcSchema.DirSizeArgs, rpcSchema.DirSizeReply](d, "DirSize", svc.DirSize)
	RegisterUnary[rpcSchema.FileCreateArgs, rpcSchema.FileCreateReply](d, "FileCreate", svc.FileCreate)
	RegisterUnary[rpcSchema.FileDeleteArgs, rpcSchema.FileDeleteReply](d, "FileDelete", svc.FileDelete)
	RegisterUnary[rpcSchema.FileUploadArgs, rpcSchema.FileUploadReply](d, "FileUpload", svc.FileUpload)
	RegisterUnary[rpcSchema.FileDownloadArgs, rpcSchema.FileDownloadReply](d, "FileDownload", svc.FileDownload)
	RegisterUnary[rpcSchema.FolderCreateArgs, rpcSchema.FolderCreateReply](d, "FolderCreate", svc.FolderCreate)
}
