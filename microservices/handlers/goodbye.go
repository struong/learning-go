package handlers

import (
	"log"
	"net/http"
)

type Goodbye struct {
	logger *log.Logger
}

func NewGoodbye(logger *log.Logger) *Goodbye {
	goodbye := new(Goodbye)
	goodbye.logger = logger

	return goodbye
}

func (g *Goodbye) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	rw.Write([]byte("Goodbye"))
}
