package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	logger *log.Logger
}

func NewHello(logger *log.Logger) *Hello {
	hello := new(Hello)
	hello.logger = logger
	return hello
}

func (h *Hello) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

	h.logger.Println("Hello World")
	data, errs := ioutil.ReadAll(request.Body)

	if errs != nil {
		http.Error(responseWriter, "Ooops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(responseWriter, "Hello %s", data)
}
