package data

import "testing"

func TestCheckValidation(t *testing.T) {
	product := new(Product)
	product.Name = "TestName"
	product.Price = 1.0
	product.SKU = "abc-abc-absd"

	err := product.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
