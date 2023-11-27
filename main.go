package main

import (
	"fmt"
	"gate-way-demo/rpc/server"
)

func main() {
	//create
	runRpcWithHttp()
	fmt.Println("end")
}

func runRpc() {
	server.Run()
}

func runRpcWithHttp() {
	server.RunHttpWrapGin()
	// server.RunWithHttp()
}
