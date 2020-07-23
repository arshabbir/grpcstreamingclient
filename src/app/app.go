package app

import (
	"log"
	"os"

	"github.com/arshabbir/grpcstreamingclient/src/client"
)

func StartApplication() {

	log.Println("Starting Client Application.......")
	os.Setenv("SERVERADDR", "localhost:9090")

	client := client.NewClient()

	client.Start()

	return
}
