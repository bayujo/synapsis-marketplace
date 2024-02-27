package cart

import (
	"context"

	"github.com/bayujo/synapsis-marketplace/internal/product"
)

type ShoppingCartUseCase struct {
	cartRepo    ShoppingCartRepository
	productRepo product.ProductRepository
}

func NewShoppingCartUseCase(cartRepo ShoppingCartRepository, productRepo product.ProductRepository) ShoppingCartUseCase {
	return ShoppingCartUseCase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (uc *ShoppingCartUseCase) AddProductToCart(ctx context.Context, userID, productID, quantity int) error {
	
	_, err := uc.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}
	
	cartID, err := uc.cartRepo.GetCartIDByUserID(ctx, userID)
	if err != nil {
		
		cartID, err = uc.cartRepo.CreateCart(ctx, userID)
		if err != nil {
			return err
		}
	}
	
	items, err := uc.cartRepo.GetCartItems(ctx, cartID)
	if err != nil {
		return err
	}
	
	for _, item := range items {
		if item.ProductID == productID {
			
			newQuantity := item.Quantity + quantity
			err := uc.cartRepo.UpdateCartItemQuantity(ctx, cartID, productID, newQuantity)
			if err != nil {
				return err
			}
			return nil
		}
	}
	
	err = uc.cartRepo.AddProductToCart(ctx, cartID, productID, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ShoppingCartUseCase) GetCartItems(ctx context.Context, userID int) ([]*CartItem, error) {
	
	cartID, err := uc.cartRepo.GetCartIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	items, err := uc.cartRepo.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, err
	}
	
	for _, item := range items {
		productInfo, err := uc.productRepo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		item.Product = productInfo
	}

	return items, nil
}

func (uc *ShoppingCartUseCase) DeleteCartItem(ctx context.Context, userID, productID int) error {
	cartID, err := uc.cartRepo.GetCartIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	err = uc.cartRepo.DeleteCartItem(ctx, cartID, productID)
	if err != nil {
		return err
	}

	return nil
}
