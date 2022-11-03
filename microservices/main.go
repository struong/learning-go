package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request) {
		log.Println("Hello World")
	})

	// create a web service, bind to every address on my machine to port 9090
	http.ListenAndServe(":9090", nil)
}
