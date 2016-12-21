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
	log.Fatal(router.Start())
}
