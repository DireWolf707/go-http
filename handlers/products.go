// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
//	Consumes:
// 	- application/json
//
//	Produces:
// 	- application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/microservices.v1/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products listProducts
// Returns a list of products
// respnses:
//	200: productResponse
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get Request")
	// getting list of products
	lp := data.GetProducts()
	// jsonify
	err := lp.ToJson(rw)
	// err handling
	if err != nil {
		http.Error(rw, "Unable to marshal", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Request")
	prod := r.Context().Value("ProductKey").(*data.Product)
	data.AddProduct(prod)
}

func (p Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Put Request")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "unable to convert id", http.StatusBadRequest)
		return
	}
	prod := r.Context().Value("ProductKey").(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}
}

// midlewares
func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJson(r.Body)
		// err handling
		if err != nil {
			http.Error(rw, "Unable to unmarshal", http.StatusBadRequest)
			return
		}
		//validating retrieved data
		err = prod.Validate()
		if err != nil {
			p.l.Println("Validation error", err)
			http.Error(
				rw,
				fmt.Sprintf("Validation error : %s", err),
				http.StatusBadRequest,
			)
			return
		}
		ctx := context.WithValue(r.Context(), "ProductKey", prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
