// Package rpc
// Date: 2024/06/25 00:15:52
// Author: Amu
// Description:
package rpc

import (
	"context"
	"os"
	"path/filepath"

	"github.com/amuluze/amvector/service/schema"
)

func (s *Service) FilesSearch(ctx context.Context, args schema.FilesSearchArgs, reply *schema.FilesSearchReply) error {
	files, err := os.ReadDir(args.Path)
	if err != nil {
		return err
	}
	data := make([]schema.FileInfo, 0)
	for _, file := range files {
		info, _ := file.Info()
		data = append(data, schema.FileInfo{
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

func (s *Service) DirSize(ctx context.Context, args schema.DirSizeArgs, reply *schema.DirSizeReply) error {
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
