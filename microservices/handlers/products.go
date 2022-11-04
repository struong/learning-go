package handlers

import (
	"log"
	"microservices/data"
	"net/http"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	product := new(Products)
	product.logger = logger
	return product
}

func (p *Products) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		p.getProducts(responseWriter, request)
		return
	}

	// handle an update
	if request.Method == http.MethodPost {
		p.addProduct(responseWriter, request)
		return
	}

	// catch all
	responseWriter.WriteHeader(http.StatusMethodNotAllowed)
}

func (products *Products) getProducts(responseWriter http.ResponseWriter, request *http.Request) {
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

func (products *Products) addProduct(responseWriter http.ResponseWriter, request *http.Request) {
	products.logger.Println("Handle POST Products")

	prod := &data.Product{}
	err := prod.FromJson(request.Body)

	if err != nil {
		products.logger.Fatal(err)
		http.Error(responseWriter, "Unable to unmarshal json", http.StatusBadRequest)
	}

	products.logger.Printf("Prod: %#v", prod)
}
