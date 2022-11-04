package main

import (
	"log"
	"microservices/handlers"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	helloHandler := handlers.NewHello(logger)
	goodbyeHandler := handlers.NewGoodbye(logger)

	serveMux := http.NewServeMux()
	// Register a handler to path /
	serveMux.Handle("/", helloHandler)
	serveMux.Handle("/goodbye", goodbyeHandler)

	server := new http.Server{}


	// create a web service, bind to every address on my machine to port 9090
	http.ListenAndServe(":9090", serveMux)
}
