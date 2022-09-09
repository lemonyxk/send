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
	"os"
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
	Name string
	Size int64
}

func sendFile(path string) {

	var f, err = os.Open(path)
	if err != nil {
		console.Exit(err)
	}

	stat, err := f.Stat()
	if err != nil {
		console.Exit(err)
	}

	var info = FileInfo{
		Name: stat.Name(),
		Size: stat.Size(),
	}

	var bts = utils.Json.Encode(info)

	var tcpSyncClient = async.NewClient[client.Conn](tcpClient)

	stream, err := tcpSyncClient.Emit("/server/fileInfo", bts)
	if err != nil {
		console.Exit(err)
	}

	if string(stream.Data) != "OK" {
		console.Exit(err)
	}

	s := make([]byte, 1024*1024*4)
	for {
		switch nr, err := f.Read(s[:]); true {
		case nr < 0:
			console.Exit(err)
		case nr == 0: // EOF
			return
		case nr > 0:
			stream, err := tcpSyncClient.Emit("/server/fileData", s[0:nr])
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
