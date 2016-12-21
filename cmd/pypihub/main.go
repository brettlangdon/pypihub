package main

import (
	"log"

	"github.com/brettlangdon/pypihub"
)

func main() {
	var config pypihub.Config
	config = pypihub.ParseConfig()

	var router *pypihub.Router
	router = pypihub.NewRouter(config)

	log.Printf("server listening on %s", config.Bind)
	log.Fatal(router.Start())
}
