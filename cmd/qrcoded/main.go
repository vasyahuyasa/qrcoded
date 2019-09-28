package main

import (
	"flag"
	"fmt"
)

const (
	defaultPort = 80
	defaultHost = "0.0.0.0"
)

func main() {
	port, host, debug := readFlags()
	run(fmt.Sprintf("%s:%d", host, port), debug)
}

func readFlags() (port int, host string, debug bool) {
	flag.IntVar(&port, "p", defaultPort, "port")
	flag.StringVar(&host, "h", defaultHost, "host")
	flag.BoolVar(&debug, "d", false, "verbose information")

	flag.Parse()
	return port, host, debug
}
