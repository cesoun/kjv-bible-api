package main

import (
	"flag"

	"github.com/cesoun/kjv-bible-api/app"
)

var (
	address = flag.String("a", "0.0.0.0", "address to bind - e.g. -a=127.0.0.1")
	port    = flag.String("p", "8080", "server port to bind - e.g. -p=8080")
)

func init() {
	flag.Parse()
}

func main() {
	a := app.App{}
	a.Init(*address, *port)
	a.Run()
}
