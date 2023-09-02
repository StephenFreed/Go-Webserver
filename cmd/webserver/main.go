package main

import (
	"flag"
	"fmt"
	"log"

	"Go-Webserver/internal/server"
)


func main() {

    // default port to 3001 | Example: ./main.go --port 8080
    var listenAddr *string = flag.String("port", ":3001", "Server Port Address")
    flag.Parse()

    var server *server.Server = server.NewServer(*listenAddr)

    fmt.Printf("Server Running On Port: %v \n", *listenAddr)
    log.Fatal(server.Start())
    fmt.Printf("Server Closed On Port: %v \n", *listenAddr)
}
