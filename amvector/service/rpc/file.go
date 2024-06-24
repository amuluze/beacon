// Package rpc
// Date: 2024/06/25 00:15:52
// Author: Amu
// Description:
package rpc

import (
	"context"
	"fmt"
	"os"

	"github.com/amuluze/amvector/service/schema"
)

func (s *Service) FilesSearch(ctx context.Context, args schema.FilesSearchArgs, reply *schema.FilesSearchReply) error {
	files, err := os.ReadDir(args.Path)
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Println(file.Name())
		fmt.Println(file.Info())
	}
	return nil
}
