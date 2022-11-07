package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products is a collection of Product
type Products []*Product

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
func (products *Products) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(products)
}

func (product *Product) FromJson(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(product)
}

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

func AddProduct(product *Product) {
	product.ID = getNextId()
	productList = append(productList, product)
}

func UpdateProduct(id int, product *Product) error {
	_, index, err := findProduct(id)
	if err != nil {
		return err
	}

	product.ID = id
	productList[index] = product
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func getNextId() int {
	last := productList[len(productList)-1]
	return last.ID + 1
}

func (p *Product) Validate() error {
	validate := validator.New()

	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {

	// sku is of format: abc-absd-dsdsd
	regex := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := regex.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
