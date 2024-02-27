package product

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	AddProduct(ctx context.Context, product *Product) error
	GetProductsByCategory(ctx context.Context, categoryID int) ([]*Product, error)
	AddProductToCart(ctx context.Context, userID, productID, quantity int) error
	GetProductByID(ctx context.Context, productID int) (*Product, error)
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) AddProduct(ctx context.Context, product *Product) error {
	_, err := r.db.Exec("INSERT INTO Products (name, description, price, quantity, category_id) VALUES ($1, $2, $3, $4, $5)",
		product.Name, product.Description, product.Price, product.Quantity, product.CategoryID)
	return err
}

func (r *productRepository) GetProductsByCategory(ctx context.Context, categoryID int) ([]*Product, error) {
	rows, err := r.db.Query("SELECT id, name, description, price, quantity, category_id FROM Products WHERE category_id = $1", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) AddProductToCart(ctx context.Context, userID, productID, quantity int) error {
	_, err := r.db.Exec("INSERT INTO Cart_Items (cart_id, product_id, quantity) VALUES ($1, $2, $3)",
		userID, productID, quantity)
	return err
}

func (r *productRepository) GetProductByID(ctx context.Context, productID int) (*Product, error) {
	var product Product
	err := r.db.QueryRow("SELECT id, name, description, price, quantity, category_id FROM Products WHERE id = $1", productID).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.CategoryID)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
