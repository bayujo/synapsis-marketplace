package product

import (
	"context"
	"errors"
)

type ProductUsecase interface {
	AddProduct(ctx context.Context, name, description string, price float64, quantity, categoryID int) error
	GetProductsByCategory(ctx context.Context, categoryID int) ([]*Product, error)
	GetProductByID(ctx context.Context, categoryID int) (*Product, error)
	AddProductToCart(ctx context.Context, userID, productID, quantity int) error
}

type productUsecase struct {
	productRepository ProductRepository
}

func NewProductUsecase(productRepository ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepository: productRepository,
	}
}

func (uc *productUsecase) AddProduct(ctx context.Context, name, description string, price float64, quantity, categoryID int) error {
	if name == "" || price <= 0 || quantity < 0 || categoryID <= 0 {
		return errors.New("invalid input")
	}

	return uc.productRepository.AddProduct(ctx, &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		CategoryID:  categoryID,
	})
}

func (uc *productUsecase) GetProductsByCategory(ctx context.Context, categoryID int) ([]*Product, error) {
	return uc.productRepository.GetProductsByCategory(ctx, categoryID)
}

func (uc *productUsecase) AddProductToCart(ctx context.Context, userID, productID, quantity int) error {
	return uc.productRepository.AddProductToCart(ctx, userID, productID, quantity)
}

func (uc *productUsecase) GetProductByID(ctx context.Context, categoryID int) (*Product, error) {
	return uc.productRepository.GetProductByID(ctx, categoryID)
}
