// Package main
// Date: 2024/7/3 18:05
// Author: Amu
// Description:
package main

import (
	"fmt"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/share"
	"io"
	"net"
	"os"
)

func main() {
	os.Remove("/tmp/rpc.sock")
	srv := server.NewServer()
	p := server.NewFileTransfer("0.0.0.0:9527", UploadFileHandler, DownloadFileHandler, 1000)
	srv.EnableFileTransfer(share.SendFileServiceName, p)

	err := srv.Serve("unix", "/tmp/rpc.sock")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func UploadFileHandler(conn net.Conn, args *share.FileTransferArgs) {
	fmt.Printf("received file name: %s, size: %d, meta: %#v\n", args.FileName, args.FileSize, args.Meta)
	data, err := io.ReadAll(conn)
	if err != nil {
		fmt.Printf("error read: %v\n", err)
		return
	}
	fmt.Printf("file content: %s\n", data)
	return
}

func DownloadFileHandler(conn net.Conn, args *share.DownloadFileArgs) {
	fmt.Printf("received file name: %s, meta: %#v\n", args.FileName, args.Meta)
	content, _ := os.Open(args.FileName)
	io.Copy(conn, content)
	conn.Close()
}
