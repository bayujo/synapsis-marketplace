package order

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
    CreateOrder(ctx context.Context, order *Order) error
    GetOrderByID(ctx context.Context, orderID int) (*Order, error)
    UpdateOrderTotalAmount(ctx context.Context, orderID int, totalAmount float64) error
    UpdateOrderStatus(ctx context.Context, orderID int, status OrderStatus) error
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = tx.QueryRow("INSERT INTO Orders (customer_id, total_amount) VALUES ($1, $2) RETURNING id", order.CustomerID, order.TotalAmount).Scan(&order.ID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO Order_Items (order_id, product_id, quantity, unit_price) VALUES ($1, $2, $3, $4)",
			order.ID, item.ProductID, item.Quantity, item.UnitPrice)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, orderID int) (*Order, error) {
    order := &Order{}
    err := r.db.Get(order, "SELECT id, customer_id, total_amount, status FROM Orders WHERE id = $1", orderID)
    if err != nil {
        return nil, err
    }

    orderItems := []*OrderItem{}
    err = r.db.Select(&orderItems, "SELECT id, order_id, product_id, quantity, unit_price FROM Order_Items WHERE order_id = $1", orderID)
    if err != nil {
        return nil, err
    }
    order.Items = orderItems

    return order, nil
}




func (r *orderRepository) UpdateOrderTotalAmount(ctx context.Context, orderID int, totalAmount float64) error {
	_, err := r.db.Exec("UPDATE Orders SET total_amount = $1 WHERE id = $2", totalAmount, orderID)
	return err
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID int, status OrderStatus) error {
    _, err := r.db.Exec("UPDATE Orders SET status = $1 WHERE id = $2", status, orderID)
    return err
}