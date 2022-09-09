/**
* @program: send
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-09 15:42
**/

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/kitty/v2"
	"github.com/lemonyxk/kitty/v2/router"
	"github.com/lemonyxk/kitty/v2/socket"
	"github.com/lemonyxk/kitty/v2/socket/tcp/server"
	"github.com/lemonyxk/utils/v3"
)

var tcpServer *server.Server

func runTcpServer() {

	var addr = ":"

	if serverAddr != "" {
		addr = serverAddr
	}

	tcpServer = kitty.NewTcpServer(addr)

	tcpServer.HeartBeatTimeout = time.Second * 5

	// tcpServer.CertFile = "example/ssl/localhost+2.pem"
	// tcpServer.KeyFile = "example/ssl/localhost+2-key.pem"

	var route = kitty.NewTcpServerRouter()

	route.Group("/server").Handler(func(handler *router.Handler[*socket.Stream[server.Conn]]) {
		handler.Route("/fileData").Handler(fileData)
		handler.Route("/fileInfo").Handler(fileInfo)
		handler.Route("/str").Handler(str)
	})

	tcpServer.OnSuccess = func() {
		console.Info("YOU ARE LISTEN ON:", tcpServer.LocalAddr().String())
		console.Info("YOU SAVE PATH IS:", savePath)
	}

	tcpServer.OnOpen = func(conn server.Conn) {}

	tcpServer.OnClose = func(conn server.Conn) {}

	tcpServer.SetRouter(route).Start()
}

var info FileInfo
var current int64 = 0
var file *os.File
var t time.Time

func fileInfo(stream *socket.Stream[server.Conn]) error {
	var err = utils.Json.Decode(stream.Data, &info)
	if err != nil {
		console.Error(err)
		return stream.Emit(stream.Event, nil)
	}

	file, err = os.OpenFile(filepath.Join(savePath, info.Name),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return stream.Emit(stream.Event, nil)
	}

	current = 0
	t = time.Now()

	console.Info("RECEIVING FILE:", info.Name, info.Size/1024/1024, "MB")

	return stream.Emit(stream.Event, []byte("OK"))
}

func fileData(stream *socket.Stream[server.Conn]) error {

	_, _ = file.Write(stream.Data)
	current += int64(len(stream.Data))

	if current == info.Size {
		_ = file.Close()
		console.Info("\nRECEIVED FILE:", info.Name, info.Size/1024/1024, "MB")
		console.Info("TIME:", float64(time.Since(t).Milliseconds())/1000, "SECONDS")
	} else {
		console.OneLine("CURRENT: %d MB TOTAL: %d MB", current/1024/1024, info.Size/1024/1024)
	}

	return stream.Emit(stream.Event, []byte("OK"))
}

func str(stream *socket.Stream[server.Conn]) error {
	console.Info(string(stream.Data))
	return nil
}
