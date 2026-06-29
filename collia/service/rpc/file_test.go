package rpc

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	rpcSchema "common/rpc/schema"
)

// newFileService 构造一个以临时目录为沙箱根的 Service，用于文件操作测试。
func newFileService(t *testing.T) (*Service, string) {
	t.Helper()
	root := t.TempDir()
	// 复用未导出的 NewService：它接受 rootDir。
	return NewService(nil, nil, root), root
}

// TestFilesSearch_RejectsOutsideSandbox 验证列举沙箱外目录被拒绝。
func TestFilesSearch_RejectsOutsideSandbox(t *testing.T) {
	svc, _ := newFileService(t)
	var reply rpcSchema.FilesSearchReply
	err := svc.FilesSearch(context.Background(), rpcSchema.FilesSearchArgs{Path: "/etc"}, &reply)
	if err == nil {
		t.Fatal("FilesSearch outside sandbox should error")
	}
}

// TestFilesSearch_ListsWithinSandbox 验证列举沙箱内目录正常工作。
func TestFilesSearch_ListsWithinSandbox(t *testing.T) {
	svc, root := newFileService(t)
	// 在沙箱内建一个文件
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("x"), 0600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	var reply rpcSchema.FilesSearchReply
	if err := svc.FilesSearch(context.Background(), rpcSchema.FilesSearchArgs{Path: root}, &reply); err != nil {
		t.Fatalf("FilesSearch inside sandbox failed: %v", err)
	}
	if len(reply.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(reply.Files))
	}
}

// TestFileDelete_RejectsOutsideSandbox 验证删除沙箱外文件被拒绝（最危险的操作）。
func TestFileDelete_RejectsOutsideSandbox(t *testing.T) {
	svc, _ := newFileService(t)
	// 在沙箱外建一个临时文件，确认删除被拒绝而非真的删掉。
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("x"), 0600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	var reply rpcSchema.FileDeleteReply
	if err := svc.FileDelete(context.Background(), rpcSchema.FileDeleteArgs{Filepath: outside}, &reply); err == nil {
		t.Fatal("FileDelete outside sandbox should error")
	}
	// 确认文件仍在（未被删）
	if _, err := os.Stat(outside); os.IsNotExist(err) {
		t.Fatal("outside file was deleted despite sandbox rejection")
	}
}

// TestFileDelete_DeletesWithinSandbox 验证删除沙箱内文件正常。
func TestFileDelete_DeletesWithinSandbox(t *testing.T) {
	svc, root := newFileService(t)
	target := filepath.Join(root, "doomed.txt")
	if err := os.WriteFile(target, []byte("x"), 0600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	var reply rpcSchema.FileDeleteReply
	if err := svc.FileDelete(context.Background(), rpcSchema.FileDeleteArgs{Filepath: target}, &reply); err != nil {
		t.Fatalf("FileDelete inside sandbox failed: %v", err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatal("file should have been deleted")
	}
}

// TestFileUpload_RejectsOutsideSandbox 验证上传到沙箱外被拒绝。
func TestFileUpload_RejectsOutsideSandbox(t *testing.T) {
	svc, _ := newFileService(t)
	var reply rpcSchema.FileUploadReply
	err := svc.FileUpload(context.Background(), rpcSchema.FileUploadArgs{
		TargetFilePath: "/etc/evil.txt",
		Data:           []byte("evil"),
	}, &reply)
	if err == nil {
		t.Fatal("FileUpload outside sandbox should error")
	}
	// 确认未写入
	if _, err := os.Stat("/etc/evil.txt"); !os.IsNotExist(err) {
		_ = os.Remove("/etc/evil.txt")
		t.Fatal("file was written outside sandbox")
	}
}

// TestFileDownload_RejectsOutsideSandbox 验证读取沙箱外文件被拒绝（防读 /etc/shadow）。
func TestFileDownload_RejectsOutsideSandbox(t *testing.T) {
	svc, _ := newFileService(t)
	var reply rpcSchema.FileDownloadReply
	err := svc.FileDownload(context.Background(), rpcSchema.FileDownloadArgs{
		SourceFilePath: "/etc/hostname",
		TargetFilePath: "", // 触发直接读取回传分支
	}, &reply)
	if err == nil {
		t.Fatal("FileDownload outside sandbox should error")
	}
}

// TestFileDownload_ReadsWithinSandbox 验证下载沙箱内文件正常。
func TestFileDownload_ReadsWithinSandbox(t *testing.T) {
	svc, root := newFileService(t)
	src := filepath.Join(root, "dl.txt")
	if err := os.WriteFile(src, []byte("content"), 0600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	var reply rpcSchema.FileDownloadReply
	if err := svc.FileDownload(context.Background(), rpcSchema.FileDownloadArgs{
		SourceFilePath: src,
		TargetFilePath: "",
	}, &reply); err != nil {
		t.Fatalf("FileDownload inside sandbox failed: %v", err)
	}
	if string(reply.Data) != "content" {
		t.Fatalf("downloaded data = %q, want content", string(reply.Data))
	}
}

// TestFileCreate_RejectsOutsideSandbox 验证在沙箱外创建文件被拒绝。
func TestFileCreate_RejectsOutsideSandbox(t *testing.T) {
	svc, _ := newFileService(t)
	var reply rpcSchema.FileCreateReply
	err := svc.FileCreate(context.Background(), rpcSchema.FileCreateArgs{
		Path:     "/etc",
		FileName: "evil.txt",
	}, &reply)
	if err == nil {
		t.Fatal("FileCreate outside sandbox should error")
	}
}

// TestFolderCreate_RejectsOutsideSandbox 验证在沙箱外创建目录被拒绝。
func TestFolderCreate_RejectsOutsideSandbox(t *testing.T) {
	svc, _ := newFileService(t)
	var reply rpcSchema.FolderCreateReply
	err := svc.FolderCreate(context.Background(), rpcSchema.FolderCreateArgs{
		Path:       "/tmp",
		FolderName: "amprobe_evil",
	}, &reply)
	if err == nil {
		t.Fatal("FolderCreate outside sandbox should error")
	}
	// 清理可能的残留
	_ = os.RemoveAll("/tmp/amprobe_evil")
}

// TestDirSize_RejectsOutsideSandbox 验证计算沙箱外目录大小被拒绝。
func TestDirSize_RejectsOutsideSandbox(t *testing.T) {
	svc, _ := newFileService(t)
	var reply rpcSchema.DirSizeReply
	err := svc.DirSize(context.Background(), rpcSchema.DirSizeArgs{Path: "/etc"}, &reply)
	if err == nil {
		t.Fatal("DirSize outside sandbox should error")
	}
}
