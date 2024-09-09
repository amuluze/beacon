// Package rpc
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package schema

type FilesSearchArgs struct {
	Path string `json:"path"`
}

type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Mode    string `json:"mode"`
	ModTime int64  `json:"mod_time"`
	IsDir   bool   `json:"is_dir"`
}
type FilesSearchReply struct {
	Files []FileInfo `json:"files"`
}

type DirSizeArgs struct {
	Path string `json:"path"`
}

type DirSizeReply struct {
	Size int64 `json:"size"`
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
	Filepath string `json:"filepath"`
}
