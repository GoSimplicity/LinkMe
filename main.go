package main

import "LinkMe/config"

func main() {
	Init()
}

func Init() {
	config.InitViper()
	cmd := InitWebServer()
	server := cmd.server
	err := cmd.consumer.Start()
	if err != nil {
		panic(err)
	}
	if er := server.Run(":9999"); er != nil {
		panic(er)
	}
}
