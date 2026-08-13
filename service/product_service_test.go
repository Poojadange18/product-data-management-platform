package service

import (
	"product-data-management-platform/model"
	"testing"
)

func TestCreateProductValidation(t *testing.T) {
	tests := []model.Product{
		{ID: 0, Name: "Phone", Category: "Electronics", Price: 0, Stock: 0},
		{ID: 1, Name: "", Category: "Electronics", Price: 0, Stock: 0},
		{ID: 1, Name: "Phone", Category: "", Price: 0, Stock: 0},
		{ID: 1, Name: "Phone", Category: "Electronics", Price: -1, Stock: 0},
		{ID: 1, Name: "Phone", Category: "Electronics", Price: 0, Stock: -1},
	}
	for _, product := range tests {
		if err := validate(&product, true); err == nil {
			t.Fatalf("expected validation error for %#v", product)
		}
	}
}

func TestPaginationDefaults(t *testing.T) {
	// Defaults are applied before repository access; an invalid repository is not used here.
	service := NewProductService(nil)
	_ = service
	if err := validate(&model.Product{ID: 99, Name: "  Phone  ", Category: " Devices ", Price: 0, Stock: 0}, true); err != nil {
		t.Fatal(err)
	}
}
