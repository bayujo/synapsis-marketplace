package cart

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type ShoppingCartRepository interface {
	CreateCart(ctx context.Context, userID int) (int, error)
	AddProductToCart(ctx context.Context, cartID, productID, quantity int) error
	GetCartItems(ctx context.Context, cartID int) ([]*CartItem, error)
	GetCartIDByUserID(ctx context.Context, userID int) (int, error)
    UpdateCartItemQuantity(ctx context.Context, cartID, productID, quantity int) error
    DeleteCartItem(ctx context.Context, cartID, productID int) error
    ClearCart(ctx context.Context, cartID int) error
}

type shoppingCartRepository struct {
	db *sqlx.DB
}

func NewShoppingCartRepository(db *sqlx.DB) ShoppingCartRepository {
	return &shoppingCartRepository{
		db: db,
	}
}

func (r *shoppingCartRepository) UpdateCartItemQuantity(ctx context.Context, cartID, productID, quantity int) error {
    _, err := r.db.Exec("UPDATE Cart_Items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3", quantity, cartID, productID)
    if err != nil {
        return err
    }
    return nil
}

// AddProductToCart implements ShoppingCartRepository.
func (r *shoppingCartRepository) AddProductToCart(ctx context.Context, cartID, productID, quantity int) error {
    // Check if the product is already in the cart
    var existingQuantity int
    err := r.db.QueryRow("SELECT quantity FROM Cart_Items WHERE cart_id = $1 AND product_id = $2", cartID, productID).Scan(&existingQuantity)
    if err == sql.ErrNoRows {
        // Product not in cart, insert new record
        _, err := r.db.Exec("INSERT INTO Cart_Items (cart_id, product_id, quantity) VALUES ($1, $2, $3)", cartID, productID, quantity)
        if err != nil {
            return err
        }
    } else if err != nil {
        return err
    } else {
        // Product already in cart, update quantity
        newQuantity := existingQuantity + quantity
        _, err := r.db.Exec("UPDATE Cart_Items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3", newQuantity, cartID, productID)
        if err != nil {
            return err
        }
    }
    return nil
}

func (r *shoppingCartRepository) CreateCart(ctx context.Context, userID int) (int, error) {
    var cartID int
    err := r.db.QueryRow("INSERT INTO Shopping_Carts (customer_id) VALUES ($1) RETURNING id", userID).Scan(&cartID)
    if err != nil {
        return 0, err
    }
    return cartID, nil
}

func (r *shoppingCartRepository) GetCartItems(ctx context.Context, cartID int) ([]*CartItem, error) {
    rows, err := r.db.Query("SELECT id, cart_id, product_id, quantity FROM Cart_Items WHERE cart_id = $1", cartID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cartItems []*CartItem
    for rows.Next() {
        var ci CartItem
        if err := rows.Scan(&ci.ID, &ci.CartID, &ci.ProductID, &ci.Quantity); err != nil {
            return nil, err
        }
        cartItems = append(cartItems, &ci)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return cartItems, nil
}

func (r *shoppingCartRepository) GetCartIDByUserID(ctx context.Context, userID int) (int, error) {
	var cartID int
	err := r.db.QueryRow("SELECT id FROM Shopping_Carts WHERE customer_id = $1", userID).Scan(&cartID)
	if err != nil {
		return 0, err
	}
	return cartID, nil
}

func (r *shoppingCartRepository) DeleteCartItem(ctx context.Context, cartID, productID int) error {
    _, err := r.db.Exec("DELETE FROM Cart_Items WHERE cart_id = $1 AND product_id = $2", cartID, productID)
    if err != nil {
        return err
    }
    return nil
}

func (r *shoppingCartRepository) ClearCart(ctx context.Context, cartID int) error {
    _, err := r.db.Exec("DELETE FROM Cart_Items WHERE cart_id = $1", cartID)
    if err != nil {
        return err
    }
    return nil
}