package order

import (
	"context"
	"fmt"

	"github.com/bayujo/synapsis-marketplace/internal/cart"
	"github.com/bayujo/synapsis-marketplace/internal/product"
	errors "github.com/bayujo/synapsis-marketplace/pkg/error"
)

type OrderUseCase interface {
    CheckoutCart(ctx context.Context, userID int) (*Order, error)
    GetOrderDetails(ctx context.Context, orderID int) (*Order, error)
    PayOrder(ctx context.Context, orderID int) error
    CancelOrder(ctx context.Context, orderID int) error
}

type orderUseCase struct {
	orderRepo   OrderRepository
	cartRepo    cart.ShoppingCartRepository
	productRepo product.ProductRepository
}

func NewOrderUseCase(orderRepo OrderRepository, cartRepo cart.ShoppingCartRepository, productRepo product.ProductRepository) OrderUseCase {
	return &orderUseCase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (uc *orderUseCase) CheckoutCart(ctx context.Context, userID int) (*Order, error) {
	
	cartID, err := uc.cartRepo.GetCartIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	items, err := uc.cartRepo.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, err
	}
	
	if len(items) == 0 {
		return nil, errors.ErrEmptyCart
	}
	
	order := &Order{
		CustomerID: userID,
		Items:      make([]*OrderItem, 0),
	}
	
	var totalAmount float64
	for _, item := range items {
		productInfo, err := uc.productRepo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		orderItem := &OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: productInfo.Price,
		}
		order.Items = append(order.Items, orderItem)

		totalAmount += float64(item.Quantity) * productInfo.Price
	}

	order.TotalAmount = totalAmount
	
	err = uc.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	
	err = uc.cartRepo.ClearCart(ctx, cartID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *orderUseCase) GetOrderDetails(ctx context.Context, orderID int) (*Order, error) {
	order, err := uc.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (uc *orderUseCase) PayOrder(ctx context.Context, orderID int) error {
	order, err := uc.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != OrderStatusPending {
		return fmt.Errorf("cannot pay order with status: %v", order.Status)
	}
	
	err = uc.orderRepo.UpdateOrderStatus(ctx, orderID, OrderStatusPaid)
	if err != nil {
		return err
	}

	return nil
}

func (uc *orderUseCase) CancelOrder(ctx context.Context, orderID int) error {
	order, err := uc.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status == OrderStatusPaid || order.Status == OrderStatusCanceled {
		return fmt.Errorf("cannot cancel order with status: %v", order.Status)
	}
	
	err = uc.orderRepo.UpdateOrderStatus(ctx, orderID, OrderStatusCanceled)
	if err != nil {
		return err
	}

	return nil
}