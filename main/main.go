package main

import (
	"flag"
	"fmt"
	"sync"
)

func main() {
	version := flag.Bool("version", false, "show im version")
	help := flag.Bool("h", false, "show help")
	port := flag.Int("p", 6890, "set port")
	flag.Parse()

	if *version {
		fmt.Println(pkg.Version)
	}
	config.port = port
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		rpc.server()
		defer wg.Done()
	}()

	go func() {
		im.server()
		defer wg.Done()
	}()
	wg.Wait()
}
