package main

import (
	"log"
	"os"

	"github.com/graph-uk/sesame/sesameServer"
)

func main() {
	sesameServer, err := sesameServer.NewSesameServer()
	if err != nil {
		log.Fatalln("Cannot start sesame server.\r\n" + err.Error())
		os.Exit(1)
	}
	go sesameServer.ExpiredRulesKiller()

	err = sesameServer.Serve()
	if err != nil {
		log.Fatalln("Cannot start sesame server.\r\n" + err.Error())
		os.Exit(1)
	}
}
