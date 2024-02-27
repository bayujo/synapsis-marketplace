package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bayujo/synapsis-marketplace/internal/cart"
	"github.com/bayujo/synapsis-marketplace/internal/handler"
	"github.com/bayujo/synapsis-marketplace/internal/middleware"
	"github.com/bayujo/synapsis-marketplace/internal/order"
	"github.com/bayujo/synapsis-marketplace/internal/product"
	"github.com/bayujo/synapsis-marketplace/internal/user"
	"github.com/bayujo/synapsis-marketplace/pkg/db/postgres"
	"github.com/bayujo/synapsis-marketplace/pkg/db/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	redisConn, err := redis.NewRedis()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer redisConn.Close()

	dbConn, err := postgres.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize auth service
	authService := middleware.NewAuthService(os.Getenv("SECRET_KEY"), os.Getenv("TOKEN_LIFESPAN"))

	// Initialize user
	userRepository := user.NewUserRepository(dbConn)
	userUsecase := user.NewUserUsecase(userRepository, authService)
	userHandler := handler.NewUserHandler(userUsecase)

	// Initialize product
	productRepository := product.NewProductRepository(dbConn)
	productUsecase := product.NewProductUsecase(productRepository)
	productHandler := handler.NewProductHandler(productUsecase, redisConn, os.Getenv("CACHE_DURATION"))

	// Initialize cart
	cartRepository := cart.NewShoppingCartRepository(dbConn)
	cartUsecase := cart.NewShoppingCartUseCase(cartRepository, productRepository)
	cartHandler := handler.NewCartHandler(cartUsecase)

	// Initialize order
	orderRepository := order.NewOrderRepository(dbConn)
	orderUsecase := order.NewOrderUseCase(orderRepository, cartRepository, productRepository)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	// Setup HTTP server
	router := mux.NewRouter()

	// User routes with middleware
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
	userRouter.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Product routes with middleware
	productRouter := router.PathPrefix("/product").Subrouter()
	productRouter.Use(middleware.AuthMiddleware)
	productRouter.HandleFunc("/store", productHandler.AddProduct).Methods("POST")
	productRouter.HandleFunc("/category/{categoryID}", productHandler.GetProductsByCategory).Methods("GET")
	productRouter.HandleFunc("/{productID}", productHandler.GetProductByID).Methods("GET")

	// Cart routes with middleware
	cartRouter := router.PathPrefix("/cart").Subrouter()
	cartRouter.Use(middleware.AuthMiddleware)
	cartRouter.HandleFunc("/add", cartHandler.AddCartItem).Methods("POST")
	cartRouter.HandleFunc("/items/{userID}", cartHandler.GetCartItems).Methods("GET")

	// Order routes with middleware
	orderRouter := router.PathPrefix("/order").Subrouter()
	orderRouter.Use(middleware.AuthMiddleware)
	orderRouter.HandleFunc("/checkout", orderHandler.CheckoutCart).Methods("POST")
	orderRouter.HandleFunc("/details/{orderID}", orderHandler.GetOrderDetails).Methods("GET")
	orderRouter.HandleFunc("/pay/{orderID}", orderHandler.PayOrder).Methods("PUT")
	orderRouter.HandleFunc("/cancel/{orderID}", orderHandler.CancelOrder).Methods("PUT")

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
