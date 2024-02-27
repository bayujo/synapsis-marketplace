package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/groupcache/singleflight"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"

	"github.com/bayujo/synapsis-marketplace/internal/product"
	envconv "github.com/bayujo/synapsis-marketplace/pkg/timeconv"
)

type ProductHandler struct {
	ProductUsecase product.ProductUsecase
	RedisClient    redis.Conn
	Group          singleflight.Group
	Cache          *cache.Cache
    CacheDuration  time.Duration
}

func NewProductHandler(productUsecase product.ProductUsecase, redisClient redis.Conn, cacheDuration string) *ProductHandler {
    cacheDurationTime := envconv.ParseDuration(cacheDuration)
	return &ProductHandler{
		ProductUsecase: productUsecase,
		RedisClient:    redisClient,
		Group:          singleflight.Group{},
		Cache:          cache.New(cacheDurationTime, cacheDurationTime),
        CacheDuration:  cacheDurationTime,
	}
}

func (h *ProductHandler) DefaultExpiration() time.Duration {
    return h.CacheDuration
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req product.Product
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.ProductUsecase.AddProduct(ctx, req.Name, req.Description, req.Price, req.Quantity, req.CategoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	categoryIDStr := vars["categoryID"]

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	cacheKey := "GetProductsByCategory" + categoryIDStr

	cachedProducts, found := h.Cache.Get(cacheKey)
	if found {
		json.NewEncoder(w).Encode(cachedProducts)
		return
	}

	result, err := h.Group.Do(cacheKey, func() (interface{}, error) {
		redisVal, err := redis.Bytes(h.RedisClient.Do("GET", cacheKey))
		if err == nil {
			var products []product.Product
			if err := json.Unmarshal(redisVal, &products); err == nil {
				h.Cache.Set(cacheKey, products, 0)
				return products, nil
			}
		}

		products, err := h.ProductUsecase.GetProductsByCategory(ctx, categoryID)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(products)
		if err != nil {
			return nil, err
		}
		_, err = h.RedisClient.Do("SET", cacheKey, data, "EX", int(h.DefaultExpiration()/time.Second))
		if err != nil {
			return nil, err
		}

		h.Cache.Set(cacheKey, products, 0)

		return products, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	productIDStr := vars["productID"]

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	cacheKey := "GetProductByID" + productIDStr

	cachedProduct, found := h.Cache.Get(cacheKey)
	if found {
		json.NewEncoder(w).Encode(cachedProduct)
		return
	}

	result, err := h.Group.Do(cacheKey, func() (interface{}, error) {
		redisVal, err := redis.Bytes(h.RedisClient.Do("GET", cacheKey))
		if err == nil {
			var product product.Product
			if err := json.Unmarshal(redisVal, &product); err == nil {
				h.Cache.Set(cacheKey, product, 0)
				return product, nil
			}
		}

		product, err := h.ProductUsecase.GetProductByID(ctx, productID)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(product)
		if err != nil {
			return nil, err
		}
		_, err = h.RedisClient.Do("SET", cacheKey, data, "EX", int(h.DefaultExpiration()/time.Second))
		if err != nil {
			return nil, err
		}

		h.Cache.Set(cacheKey, product, 0)

		return product, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}


func (h *ProductHandler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		CartID    int `json:"cart_id"`
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.ProductUsecase.AddProductToCart(ctx, req.CartID, req.ProductID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
