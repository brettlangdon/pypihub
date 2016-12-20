package main

import (
	"fmt"

	"github.com/brettlangdon/pypihub"
)

func main() {
	var config pypihub.Config
	config = pypihub.ParseConfig()
	fmt.Printf("%#v\r\n", config)

	var server *pypihub.Server
	server = pypihub.NewServer(config)
	panic(server.ListenAndServe())
}
