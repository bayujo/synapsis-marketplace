package cart

import "github.com/bayujo/synapsis-marketplace/internal/product"

type ShoppingCart struct {
    ID         int         `json:"id"`
    CustomerID int         `json:"customer_id"`
    Items      []*CartItem `json:"items"`
}

type CartItem struct {
	ID        int              `json:"id"`
	CartID    int              `json:"cart_id"`
	ProductID int              `json:"product_id"`
	Quantity  int              `json:"quantity"`
	Product   *product.Product `json:"product"`
}