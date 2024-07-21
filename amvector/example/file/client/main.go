// Package main
// Date: 2024/7/3 17:50
// Author: Amu
// Description:
package main

import (
	"context"
	"log"
	"os"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
)

func main() {
	d, _ := client.NewPeer2PeerDiscovery("unix@/tmp/rpc.sock", "")
	//xclient := client.NewXClient(share.SendFileServiceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	//defer xclient.Close()

	//err := xclient.SendFile(context.Background(), "/Users/amu/Desktop/github/amprobe/amprobe/amvector/example/file/client/abc.txt", 0, map[string]string{"path": "/tmp"})
	//if err != nil {
	//	panic(err)
	//}
	//log.Println("send ok")

	// 写入文件 的 io.Writer
	// xclient := client.NewOneClient(client.Failtry, client.RandomSelect, d, client.DefaultOption)
	xclient := client.NewXClient(share.SendFileServiceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	file, _ := os.Create("/Users/amu/Desktop/github/amprobe/amprobe/amvector/example/file/client/ddd.txt")
	defer file.Close()
	err := xclient.DownloadFile(context.Background(), "/Users/amu/Desktop/github/amprobe/amprobe/amvector/example/file/client/abc.txt", file, map[string]string{})
	if err != nil {
		panic(err)
	}
	log.Println("download ok")
}
