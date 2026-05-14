package main

import (
	"log"

	"github.com/sagar-rathod-tech53/provenloop/routes"
)

func main() {

	router := routes.SetupServer()

	err := router.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
