// Package rpc
// Date: 2024/06/25 00:15:52
// Author: Amu
// Description:
package rpc

import (
	"amvector/pkg/utils"
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
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.FileMode(0755))
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (s *Service) FolderCreate(ctx context.Context, args rpcSchema.FolderCreateArgs, reply *rpcSchema.FolderCreateReply) error {
	folderPath := filepath.Join(args.Path, args.FolderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return os.Mkdir(folderPath, os.FileMode(0755))
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
	return os.Rename(args.SourceFilePath, args.TargetFilePath)
}

func (s *Service) FileDownload(ctx context.Context, args rpcSchema.FileDownloadArgs, reply *rpcSchema.FileDownloadReply) error {
	_, err := utils.CopyFile(args.SourceFilePath, args.TargetFilePath)
	reply.Filepath = args.TargetFilePath
	return err
}
