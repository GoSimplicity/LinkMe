package main

import "LinkMe/config"

func main() {
	config.InitViper()
	server := InitWebServer()
	if err := server.Run(":9999"); err != nil {
		panic(err)
	}
}
