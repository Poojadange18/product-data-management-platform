package service

import (
	"errors"
	"product-data-management-platform/model"
	"product-data-management-platform/repository"
	"strings"
)

type ProductService struct{ Repository *repository.ProductRepository }

func NewProductService(r *repository.ProductRepository) *ProductService {
	return &ProductService{Repository: r}
}
func validate(p *model.Product, requireID bool) error {
	p.Name = strings.TrimSpace(p.Name)
	p.Category = strings.TrimSpace(p.Category)
	if requireID && p.ID <= 0 {
		return errors.New("product ID must be positive")
	}
	if p.Name == "" {
		return errors.New("product name is required")
	}
	if p.Category == "" {
		return errors.New("product category is required")
	}
	if p.Price < 0 {
		return errors.New("product price cannot be negative")
	}
	if p.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	return nil
}
func (s *ProductService) CreateProduct(p *model.Product) error {
	if err := validate(p, true); err != nil {
		return err
	}
	return s.Repository.CreateProduct(p)
}
func (s *ProductService) GetProducts(q model.ProductQuery) (model.PaginatedProducts, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 20
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	return s.Repository.GetProducts(q)
}
func (s *ProductService) GetProductByID(id int) (*model.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}
	return s.Repository.GetProductByID(id)
}
func (s *ProductService) UpdateProduct(id int, p *model.Product) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}
	if err := validate(p, false); err != nil {
		return err
	}
	return s.Repository.UpdateProduct(id, p)
}
func (s *ProductService) DeleteProduct(id int) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}
	return s.Repository.DeleteProduct(id)
}
