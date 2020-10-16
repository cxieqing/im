package main

import (
	"flag"
	"fmt"
	"im/pkg/config"
	"im/server/im"
	"im/server/rpc"
)

func main() {
	version := flag.Bool("version", false, "show im version")
	addr := flag.String("addr", ":8080", "http service address")
	//help := flag.Bool("h", false, "show help")
	port := flag.Int("p", 6890, "set port")
	flag.Parse()

	if *version {
		fmt.Println(im.Version)
	}
	config := config.NewConfig()
	config.Port = *port

	var wg sync.WaitGroup
	wg.Add(3)
	go func(){
		defer wg.Done()
		im.CreateImServer()
	}
	go func(){
		defer wg.Done()
		rpc.CreateRpcServer()
	} 
	wg.Wait()
}
