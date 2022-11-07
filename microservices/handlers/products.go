package handlers

import (
	"context"
	"fmt"
	"log"
	"microservices/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	logger *log.Logger
}

type KeyProduct struct {
}

func NewProducts(logger *log.Logger) *Products {
	product := new(Products)
	product.logger = logger
	return product
}

func (products *Products) GetProducts(responseWriter http.ResponseWriter, request *http.Request) {
	products.logger.Println("Handle GET Products")

	listProducts := data.GetProducts()
	err := listProducts.ToJSON(responseWriter)

	// data, err := json.Marshal(listProducts)

	if err != nil {
		products.logger.Fatal(err)
		http.Error(responseWriter, "Unable to marshal json", http.StatusInternalServerError)
	}

	// responseWriter.Write(data)
}

func (products *Products) AddProduct(responseWriter http.ResponseWriter, request *http.Request) {
	products.logger.Println("Handle POST Products")

	prod := request.Context().Value(KeyProduct{}).(data.Product)

	products.logger.Printf("Prod: %#v", prod)
	data.AddProduct(&prod)
}

func (products *Products) UpdateProducts(responseWriter http.ResponseWriter, request *http.Request) {
	// extract id
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(responseWriter, "Unable to convert id to int", http.StatusBadRequest)
	}

	products.logger.Printf("Handle PUT Products, %v", id)

	prod := request.Context().Value(KeyProduct{}).(data.Product)

	// products.logger.Printf("Prod: %#v", prod)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(responseWriter, data.ErrProductNotFound.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(responseWriter, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (products *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		prod := data.Product{}
		err := prod.FromJson(request.Body)

		if err != nil {
			products.logger.Println(err)
			http.Error(responseWriter, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			products.logger.Println(err)
			http.Error(
				responseWriter,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(request.Context(), KeyProduct{}, prod)
		newRequest := request.WithContext(ctx)
		next.ServeHTTP(responseWriter, newRequest)
	})
}
