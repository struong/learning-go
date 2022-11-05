package handlers

import (
	"log"
	"microservices/data"
	"net/http"
	"regexp"
	"strconv"
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

	if request.Method == http.MethodPut {
		// expect the id in the URI
		regex := regexp.MustCompile(`/([0-9]+)`)
		groups := regex.FindAllStringSubmatch(request.URL.Path, -1)

		if len(groups) != 1 {
			p.logger.Println("Invalid URI more than one id")
			http.Error(responseWriter, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(groups[0]) != 2 {
			p.logger.Println("Invalid URI more than one capture group")
			http.Error(responseWriter, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := groups[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.logger.Println("Invalid URI unable to convert to number")
			http.Error(responseWriter, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, responseWriter, request)
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
		return
	}

	products.logger.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (products *Products) updateProducts(id int, responseWriter http.ResponseWriter, request *http.Request) {
	products.logger.Println("Handle PUT Products")

	prod := &data.Product{}
	err := prod.FromJson(request.Body)

	if err != nil {
		products.logger.Fatal(err)
		http.Error(responseWriter, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	products.logger.Printf("Prod: %#v", prod)
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(responseWriter, data.ErrProductNotFound.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(responseWriter, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
