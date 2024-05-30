package main

import (
	"LinkMe/config"
)

func main() {
	config.InitViper()
	cmd := InitWebServer()

	for _, c := range cmd.consumer {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := cmd.server
	if err := server.Run(":9999"); err != nil {
		panic(err)
	}
}
