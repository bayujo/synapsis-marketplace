## Synapsis Marketplace Backend Documentation

### Project Overview
Synapsis Marketplace is an online store application backend, focusing on providing RESTful APIs for managing products, carts, orders, and user authentication.

### Entity Relationship Diagram
![Entity Relationship Diagram](https://raw.githubusercontent.com/bayujo/synapsis-marketplace/main/erd.png)

### Features
- **User Management:**
  - Register: `POST /user/register`
  - Login: `POST /user/login`
- **Product Management:**
  - Add Product: `POST /product/store`
  - Get Products by Category: `GET /product/category/{categoryID}`
  - Get Product by ID: `GET /product/{productID}`
- **Cart Management:**
  - Add Item to Cart: `POST /cart/add`
  - Get Cart Items: `GET /cart/items/{userID}`
- **Order Management:**
  - Checkout Cart: `POST /order/checkout`
  - Get Order Details: `GET /order/details/{orderID}`
  - Pay for Order: `PUT /order/pay/{orderID}`
  - Cancel Order: `PUT /order/cancel/{orderID}`

### Dependencies
- Mux: [https://github.com/gorilla/mux](https://github.com/gorilla/mux)
- JWT Authentication: [https://github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)
- Password Encryption: [https://pkg.go.dev/golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- Postgres: [https://github.com/lib/pq](https://github.com/lib/pq)
- sqlx: [https://github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx)
- Redis Caching: [https://github.com/gomodule/redigo](https://github.com/gomodule/redigo)
- In-Memory Caching: [https://pkg.go.dev/github.com/patrickmn/go-cache](https://pkg.go.dev/github.com/patrickmn/go-cache)
- Singleflight: For duplicate function call suppression mechanism
- Database Indexing: Implemented to improve performance

### Docker Setup
Dockerfile and docker-compose files are provided for containerization, build and run the image:

    docker-compose up --build

## API Testing
Open [API Documentation](https://documenter.getpostman.com/view/20497104/2sA2rDy1nY)
Import to your Postman account or export to JSON format.

## Deployment
The project is deployed to GCP platforms.