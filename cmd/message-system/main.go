package main

import (
	"log"
	"os"

	"github.com/multiplay/internal/hub"
)

const (
	PORT = ":80"
	TYPE = "tcp"
)

func main() {

	// Try server startup
	srv, err := hub.NewHub(PORT, TYPE)

	if err != nil {
		log.Println("Error Building Server Config: ", err.Error())
		os.Exit(1)
	}

	err = srv.ServerConn()

	if err != nil {
		log.Println("Error Starting Server: ", err.Error())
		os.Exit(1)
	}

}
