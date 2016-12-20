package main

import (
	"fmt"

	"github.com/brettlangdon/pypihub"
)

func main() {
	var config pypihub.Config
	config = pypihub.ParseConfig()
	fmt.Printf("%#v\r\n", config)
}
