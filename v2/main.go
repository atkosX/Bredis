package main

import (
	"flag"
	"log"
	"main/config"
	"main/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Bredis host server")
	flag.IntVar(&config.Port, "port", 6379, "Bredis port")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Bredis server baking up on ", config.Host, ":", config.Port)
	server.RunAsyncTCPServer()
}
