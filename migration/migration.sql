CREATE TABLE Categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE Products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    category_id INT NOT NULL,
    FOREIGN KEY (category_id) REFERENCES Categories(id)
);

CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE Shopping_Carts (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES Users(id)
);

CREATE TABLE Cart_Items (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (cart_id) REFERENCES Shopping_Carts(id),
    FOREIGN KEY (product_id) REFERENCES Products(id)
);

CREATE TABLE Orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10, 2) NOT NULL,
    status INT DEFAULT 0 NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES Users(id)
);

CREATE TABLE Order_Items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES Orders(id),
    FOREIGN KEY (product_id) REFERENCES Products(id)
);

CREATE INDEX idx_users_username ON Users (username);
CREATE INDEX idx_users_email ON Users (email);

CREATE INDEX idx_products_category_id ON Products (category_id);

CREATE INDEX idx_orders_id ON Orders (id);
CREATE INDEX idx_orders_customer_id ON Orders (customer_id);

CREATE INDEX idx_order_items_order_id ON Order_Items (order_id);

CREATE INDEX idx_cart_items_cart_product ON Cart_Items (cart_id, product_id);
CREATE INDEX idx_cart_items_cart_id ON Cart_Items (cart_id);

CREATE INDEX idx_shopping_carts_customer_id ON Shopping_Carts (customer_id);

INSERT INTO Categories (name, description) VALUES
('Electronics', 'Electronic devices and accessories'),
('Clothing', 'Various types of clothing items'),
('Books', 'Books on various topics'),
('Home & Kitchen', 'Items for home and kitchen use');

INSERT INTO Products (name, description, price, quantity, category_id) VALUES
('Smartphone', 'High-end smartphone with advanced features', 799.99, 100, 1),
('Laptop', 'Powerful laptop for professional use', 1299.99, 50, 1),
('Headphones', 'Wireless headphones with noise cancellation', 199.99, 80, 1),
('T-Shirt', 'Casual cotton t-shirt', 19.99, 200, 2),
('Jeans', 'Blue denim jeans', 39.99, 150, 2),
('Dress Shirt', 'Formal dress shirt for men', 49.99, 100, 2),
('Python Programming', 'Book on Python programming language', 29.99, 75, 3),
('Data Science Handbook', 'Comprehensive guide to data science', 39.99, 50, 3),
('Fiction Novel', 'Bestselling fiction novel', 14.99, 120, 3),
('Cookware Set', 'Set of pots and pans for cooking', 99.99, 80, 4),
('Blender', 'High-speed blender for smoothies and soups', 59.99, 60, 4),
('Coffee Maker', 'Automatic coffee maker with built-in grinder', 129.99, 40, 4);