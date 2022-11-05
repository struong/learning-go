package main

import (
	"context"
	"log"
	"microservices/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "product-api: ", log.LstdFlags)

	// create the handlers
	productsHandler := handlers.NewProducts(logger)

	// create a new serve mux and register the handlers
	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productsHandler.GetProducts)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", productsHandler.UpdateProducts)
	putRouter.Use(productsHandler.MiddlewareProductValidation)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productsHandler.AddProduct)
	postRouter.Use(productsHandler.MiddlewareProductValidation)

	server := new(http.Server)
	server.Addr = ":9090"
	server.Handler = serveMux
	server.IdleTimeout = 120 * time.Second
	server.ReadTimeout = 1 * time.Second
	server.WriteTimeout = 1 * time.Second

	// Do not block http server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)

	// broadcast a message on the signalChannel when we receive these os signals
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// Reading from a channel will block until there is a message to consume
	sig := <-signalChannel
	logger.Println("Recieved terminate, graceful shutdown", sig)

	// 30 seconds to gracefully shutdown
	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)

	server.Shutdown(timeoutContext)
}
