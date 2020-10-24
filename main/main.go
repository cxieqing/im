package main

import (
	"flag"
	"fmt"
	"im/pkg/config"
	"im/server/im"
	"im/server/rpc"
	"sync"
)

func main() {
	version := flag.Bool("version", false, "show im version")
	//help := flag.Bool("h", false, "show help")
	port := flag.Int("ip", 0, "set im port")
	flag.Parse()

	if *version {
		fmt.Println(im.Version)
	}
	config := config.NewConfig()
	if *port > 0 {
		config.ImPort = *port
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		im.CreateImServer()
	}()
	go func() {
		defer wg.Done()
		rpc.CreateRpcServer()
	}()
	wg.Wait()
}
