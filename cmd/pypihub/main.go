package main

import "github.com/brettlangdon/pypihub"

func main() {
	var config pypihub.Config
	config = pypihub.ParseConfig()

	var server *pypihub.Server
	server = pypihub.NewServer(config)
	panic(server.ListenAndServe())
}
