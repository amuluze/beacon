// Package schema
// Date: 2024/06/25 00:18:26
// Author: Amu
// Description:
package schema

type FilesSearchArgs struct {
	Path string
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    string
	ModTime int64
	IsDir   bool
}
type FilesSearchReply struct {
	Files []FileInfo
}

type DirSizeArgs struct {
	Path string
}

type DirSizeReply struct {
	Size int64
}
