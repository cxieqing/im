package main

import (
	"flag"
	"fmt"
	"im/pkg/config"
	"im/server"
	"net/http"
)

func main() {
	version := flag.Bool("version", false, "show im version")
	addr := flag.String("addr", ":80", "http service address")
	//help := flag.Bool("h", false, "show help")
	port := flag.Int("p", 6890, "set port")
	flag.Parse()

	if *version {
		fmt.Println(server.Version)
	}
	config := config.NewConfig()
	config.Port = *port
	http.HandleFunc("/im", server.ImServer)
	http.HandleFunc("/user", server.UserRpc)
	http.ListenAndServe(*addr, nil)
}
