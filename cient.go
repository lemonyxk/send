/**
* @program: send
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-09 15:43
**/

package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/kitty/v2"
	"github.com/lemonyxk/kitty/v2/socket/async"
	"github.com/lemonyxk/kitty/v2/socket/tcp/client"
	"github.com/lemonyxk/utils/v3"
)

var tcpClient *client.Client

func runTcpClient(addr string) chan struct{} {

	var ch = make(chan struct{}, 1)

	tcpClient = kitty.NewTcpClient(addr)

	// tcpClient.HeartBeatTimeout = time.Second * 3
	tcpClient.HeartBeatInterval = time.Second * 1
	tcpClient.WriteBufferSize = 1024 * 1024 * 4
	tcpClient.ReadBufferSize = 1024 * 1024 * 4

	var route = kitty.NewTcpClientRouter()

	// tcpClient.CertFile = "example/ssl/localhost+2.pem"
	// tcpClient.KeyFile = "example/ssl/localhost+2-key.pem"

	// make sure the event run only once
	// because when the client reconnect, the event will be run again and chan will be blocked.
	tcpClient.OnSuccess = func() {
		ch <- struct{}{}
	}

	tcpClient.OnOpen = func(conn client.Conn) {}
	tcpClient.OnClose = func(conn client.Conn) {}

	go tcpClient.SetRouter(route).Connect()

	return ch
}

type FileInfo struct {
	Name   string
	Size   int64
	Prefix string
}

func sendFile(path string) {
	stat, err := os.Stat(path)
	if err != nil {
		console.Exit(err)
	}

	if !stat.IsDir() {
		var dir = filepath.Dir(path)
		doFile(dir, dir, stat)
		return
	}

	var fn func(p string)

	fn = func(p string) {
		files, err := os.ReadDir(p)
		if err != nil {
			console.Error(err)
			return
		}

		for i := 0; i < len(files); i++ {
			var fullPath = filepath.Join(p, files[i].Name())
			if files[i].IsDir() {
				fn(fullPath)
				continue
			}

			info, err := files[i].Info()
			if err != nil {
				console.Error(err)
				continue
			}

			doFile(p, path, info)
		}
	}

	fn(path)
}

var buf = make([]byte, 1024*1024*4)

func doFile(fullPath string, rPath string, stat fs.FileInfo) {

	console.Info(filepath.Join(fullPath, stat.Name()))

	f, err := os.Open(filepath.Join(fullPath, stat.Name()))
	if err != nil {
		console.Exit(err)
	}

	defer func() { _ = f.Close() }()

	var info = FileInfo{
		Name:   stat.Name(),
		Size:   stat.Size(),
		Prefix: strings.ReplaceAll(fullPath, rPath, ""),
	}

	var bts = utils.Json.Encode(info)

	var tcpSyncClient = async.NewClient[client.Conn](tcpClient)

	stream, err := tcpSyncClient.Emit("/server/fileInfo", bts)
	if err != nil {
		console.Exit(err)
	}

	if string(stream.Data) != "OK" {
		console.Exit(string(stream.Data))
	}

	for {
		switch nr, err := f.Read(buf[:]); true {
		case nr < 0:
			console.Exit(err)
		case nr == 0: // EOF
			return
		case nr > 0:
			stream, err := tcpSyncClient.Emit("/server/fileData", buf[0:nr])
			if err != nil {
				console.Error(err)
			}
			if string(stream.Data) != "OK" {
				console.Exit(err)
			}
		}
	}
}

func sendStr(str string) {
	var err = tcpClient.Emit("/server/str", []byte(str))
	console.Info("send", str)
	if err != nil {
		console.Error(err)
	}
}
