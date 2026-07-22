// Package rpc
// Date: 2024/06/25 00:15:52
// Author: Amu
// Description:
package rpc

import (
	"collia/pkg/utils"
	"context"
	"os"
	"path/filepath"

	rpcSchema "common/rpc/schema"
)

// FilesSearch 列出指定目录下的文件。path 必须位于沙箱根目录内，
// 防止越界列举系统目录。
func (s *Service) FilesSearch(ctx context.Context, args rpcSchema.FilesSearchArgs, reply *rpcSchema.FilesSearchReply) error {
	if _, err := utils.SanitizePathWithin(args.Path, s.rootDir); err != nil {
		return err
	}
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

// DirSize 计算目录大小。path 必须位于沙箱内。
func (s *Service) DirSize(ctx context.Context, args rpcSchema.DirSizeArgs, reply *rpcSchema.DirSizeReply) error {
	if _, err := utils.SanitizePathWithin(args.Path, s.rootDir); err != nil {
		return err
	}
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

// FileCreate 创建空文件。路径必须位于沙箱内。
func (s *Service) FileCreate(ctx context.Context, args rpcSchema.FileCreateArgs, reply *rpcSchema.FileCreateReply) error {
	filePath := filepath.Join(args.Path, args.FileName)
	if _, err := utils.SanitizePathWithin(filePath, s.rootDir); err != nil {
		return err
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600) //#nosec G304 -- path sanitized by utils.SanitizePathWithin
		if err != nil {
			return err
		}
		return f.Close()
	}
	return nil
}

// FolderCreate 创建目录。路径必须位于沙箱内。
func (s *Service) FolderCreate(ctx context.Context, args rpcSchema.FolderCreateArgs, reply *rpcSchema.FolderCreateReply) error {
	folderPath := filepath.Join(args.Path, args.FolderName)
	if _, err := utils.SanitizePathWithin(folderPath, s.rootDir); err != nil {
		return err
	}
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return os.Mkdir(folderPath, 0750)
	}
	return nil
}

// FileDelete 删除文件或目录。filepath 必须位于沙箱内，
// 防止越界删除系统文件（原实现可 rm -rf 任意路径）。
func (s *Service) FileDelete(ctx context.Context, args rpcSchema.FileDeleteArgs, reply *rpcSchema.FileDeleteReply) error {
	if _, err := utils.SanitizePathWithin(args.Filepath, s.rootDir); err != nil {
		return err
	}
	if info, err := os.Stat(args.Filepath); err != nil {
		return err
	} else if info.IsDir() {
		return os.RemoveAll(args.Filepath)
	} else {
		return os.Remove(args.Filepath)
	}
}

// FileDownload 下载文件。source 必须位于沙箱内，防止越界读取
// /etc/shadow、私钥等敏感文件；target（服务端拷贝）同样校验。
func (s *Service) FileDownload(ctx context.Context, args rpcSchema.FileDownloadArgs, reply *rpcSchema.FileDownloadReply) error {
	if _, err := utils.SanitizePathWithin(args.SourceFilePath, s.rootDir); err != nil {
		return err
	}
	if args.TargetFilePath == "" {
		data, err := os.ReadFile(args.SourceFilePath) //#nosec G304 -- path sanitized
		if err != nil {
			return err
		}
		reply.Filepath = args.SourceFilePath
		reply.FileName = filepath.Base(args.SourceFilePath)
		reply.Data = data
		return nil
	}
	// CopyFile 内部对 src/dst 已做 SanitizePath 校验；这里额外约束沙箱根目录。
	if _, err := utils.SanitizePathWithin(args.TargetFilePath, s.rootDir); err != nil {
		return err
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
	RegisterUnary[rpcSchema.FileDownloadArgs, rpcSchema.FileDownloadReply](d, "FileDownload", svc.FileDownload)
	RegisterUnary[rpcSchema.FolderCreateArgs, rpcSchema.FolderCreateReply](d, "FolderCreate", svc.FolderCreate)
}
