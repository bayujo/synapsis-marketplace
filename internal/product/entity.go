package product

type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Quantity    int     `json:"quantity"`
    CategoryID  int     `json:"category_id"`
}
