package main

import (
	"log"

	"github.com/Wheeeel/pushen-server/api"
)

func main() {
	server, err := api.New(":8080")
	if err != nil {
		log.Fatal(err)
	}
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
