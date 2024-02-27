package order

type Order struct {
    ID          int          `json:"id" db:"id"`
    CustomerID  int          `json:"customer_id" db:"customer_id"`
    Items       []*OrderItem `json:"items"`
    TotalAmount float64      `json:"total_amount" db:"total_amount"`
    Status      OrderStatus  `json:"status" db:"status"`
}

type OrderItem struct {
    ID        int     `json:"id" db:"id"`
    OrderID   int     `json:"order_id" db:"order_id"`
    ProductID int     `json:"product_id" db:"product_id"`
    Quantity  int     `json:"quantity" db:"quantity"`
    UnitPrice float64 `json:"unit_price" db:"unit_price"`
}

type OrderStatus int

const (
	OrderStatusPending OrderStatus = iota
	OrderStatusPaid
	OrderStatusCanceled
)
