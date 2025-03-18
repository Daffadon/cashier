package main

import (
	"tiga-putra-cashier-be/cmd"
	"tiga-putra-cashier-be/di"
)

func main() {
	container := di.BuildContainer()
	serverReady := make(chan bool)
	server := &cmd.Server{
		Container:   container,
		ServerReady: serverReady,
	}
	server.Start()
	<-serverReady
}
