/**
* @program: send
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-09 15:40
**/

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/utils/v3"
)

func init() {
	console.SetFlags(0)
	console.Colorful(false)

	if HasArgs("--debug", os.Args) {
		console.SetFlags(console.FILE)
	}
}

var (
	savePath   = ""
	serverAddr = ""
)

func main() {

	if len(os.Args) < 2 {
		console.Exit(help())
	}

	switch os.Args[1] {
	case "server":
		startServer()
	case "client":
		startClient()
	default:
		console.Exit(help())
	}

}

func startServer() {

	var pwd, err = os.Getwd()
	if err != nil {
		console.Exit(err)
	}

	absPwd, err := filepath.Abs(pwd)
	if err != nil {
		console.Exit(err)
	}

	savePath = absPwd

	if HasArgs("--path", os.Args) {
		var p = GetArgs([]string{"--path"}, os.Args)
		if p == "" {
			console.Exit("save path is empty")
		}
		var absP, err = filepath.Abs(p)
		if err != nil {
			console.Exit(err)
		}
		savePath = absP
	}

	if HasArgs("--addr", os.Args) {
		var p = GetArgs([]string{"--addr"}, os.Args)
		if p == "" {
			console.Exit("addr is empty")
		}
		serverAddr = p
	}

	go runTcpServer()

	utils.Signal.ListenKill().Done(func(sig os.Signal) {
		console.Exit(sig)
	})
}

func startClient() {

	var addr = GetArgs([]string{"--addr"}, os.Args)
	if addr == "" {
		console.Exit("addr is empty")
	}

	<-runTcpClient(addr)

	if HasArgs("--file", os.Args) || HasArgs("-f", os.Args) {
		var p = GetArgs([]string{"--file", "-f"}, os.Args)
		var absP, err = filepath.Abs(p)
		if err != nil {
			console.Exit(err)
		}
		sendFile(absP)
		return
	}

	if HasArgs("--string", os.Args) || HasArgs("-s", os.Args) {
		var p = GetArgs([]string{"--string", "-s"}, os.Args)
		if p == "" {
			console.Exit("string is empty")
		}
		sendStr(p)
		return
	}

	time.Sleep(time.Second)
}
