package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"product-data-management-platform/model"
)

type ProductRepository struct{ DB *sql.DB }

func NewProductRepository(db *sql.DB) *ProductRepository { return &ProductRepository{DB: db} }

func (r *ProductRepository) CreateProduct(p *model.Product) error {
	return r.DB.QueryRow(`INSERT INTO products (id, name, category, price, stock) VALUES ($1, $2, $3, $4, $5) RETURNING id`, p.ID, p.Name, p.Category, p.Price, p.Stock).Scan(&p.ID)
}

func (r *ProductRepository) GetProducts(q model.ProductQuery) (model.PaginatedProducts, error) {
	where, args := []string{"1=1"}, []interface{}{}
	if q.Search != "" {
		args = append(args, "%"+q.Search+"%")
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR category ILIKE $%d)", len(args), len(args)))
	}
	if q.Category != "" {
		args = append(args, q.Category)
		where = append(where, fmt.Sprintf("category = $%d", len(args)))
	}
	switch q.Stock {
	case "in-stock":
		where = append(where, "stock >= 10")
	case "low-stock":
		where = append(where, "stock BETWEEN 1 AND 9")
	case "out-of-stock":
		where = append(where, "stock = 0")
	}
	filter := strings.Join(where, " AND ")
	result := model.PaginatedProducts{Page: q.Page, Limit: q.Limit}
	if err := r.DB.QueryRow("SELECT COUNT(*) FROM products WHERE "+filter, args...).Scan(&result.Total); err != nil {
		return result, err
	}
	result.TotalPages = (result.Total + q.Limit - 1) / q.Limit
	allowedSort := map[string]string{"name": "name", "price": "price", "stock": "stock", "id": "id"}
	sort := allowedSort[q.Sort]
	if sort == "" {
		sort = "id"
	}
	order := "ASC"
	if strings.EqualFold(q.Order, "desc") {
		order = "DESC"
	}
	args = append(args, q.Limit, (q.Page-1)*q.Limit)
	query := fmt.Sprintf("SELECT id, name, category, price, stock FROM products WHERE %s ORDER BY %s %s, id ASC LIMIT $%d OFFSET $%d", filter, sort, order, len(args)-1, len(args))
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	result.Data = []model.Product{}
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock); err != nil {
			return result, err
		}
		result.Data = append(result.Data, p)
	}
	return result, rows.Err()
}
func (r *ProductRepository) GetProductByID(id int) (*model.Product, error) {
	p := &model.Product{}
	err := r.DB.QueryRow(`SELECT id,name,category,price,stock FROM products WHERE id=$1`, id).Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (r *ProductRepository) UpdateProduct(id int, p *model.Product) error {
	result, err := r.DB.Exec(`UPDATE products SET name=$1,category=$2,price=$3,stock=$4,updated_at=CURRENT_TIMESTAMP WHERE id=$5`, p.Name, p.Category, p.Price, p.Stock, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	p.ID = id
	return nil
}
func (r *ProductRepository) DeleteProduct(id int) error {
	result, err := r.DB.Exec(`DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
