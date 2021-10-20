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

// CLI Arguments to specify how many clients to run for easier testing?
func main() {

	// Try server startup
	err := hub.ServerConn(PORT, TYPE)

	if err != nil {
		log.Println("Error Starting Server: ", err.Error())
		os.Exit(1)
	}

}
