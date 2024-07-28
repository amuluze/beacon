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

type FileCreateArgs struct {
	Path     string `json:"path" validate:"required"`
	FileName string `json:"file_name" validate:"required"`
}

type FileCreateReply struct{}

type FolderCreateArgs struct {
	Path       string `json:"path" validate:"required"`
	FolderName string `json:"folder_name" validate:"required"`
}

type FolderCreateReply struct{}

type FileDeleteArgs struct {
	Filepath string `json:"filepath" validate:"required"`
}

type FileDeleteReply struct{}

type FileUploadArgs struct {
	SourceFilePath string `json:"source_file_path" validate:"required"`
	TargetFilePath string `json:"target_file_path" validate:"required"`
}

type FileUploadReply struct{}

type FileDownloadArgs struct {
	SourceFilePath string `json:"source_file_path" validate:"required"`
	TargetFilePath string `json:"target_file_path" validate:"required"`
}

type FileDownloadReply struct {
	Filepath string
}
