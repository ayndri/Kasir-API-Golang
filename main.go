package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int 		`json:"id"`
	Name  string 	`json:"name"`
	Price int 		`json:"price"`
	Stock int 		`json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Product 1", Price: 1000, Stock: 10},
	{ID: 2, Name: "Product 2", Price: 2000, Stock: 20},
	{ID: 3, Name: "Product 3", Price: 3000, Stock: 30},
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, p := range products {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func updateProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i := range products {
		if products[i].ID == id {
			updatedProduct.ID = id
			products[i] = updatedProduct
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}

	// for i, p := range products {
	// 	if p.ID == id {
	// 		products[i] = updatedProduct
	// 		w.Header().Set("Content-Type", "application/json")
	// 		json.NewEncoder(w).Encode(updatedProduct)
	// 		return
	// 	}
	// }

	http.Error(w, "Product not found", http.StatusNotFound)
}

func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "success delete",
			})
			return
		}
	}

	// for i := range products {
	// 	if products[i].ID == id {
	// 		products = append(products[:i], products[i+1:]...)
	// 		w.WriteHeader(http.StatusNoContent)
	// 		return
	// 	}
	// }

	http.Error(w, "Product not found", http.StatusNotFound)
}

func main() {
	// GET localhost:8080/api/product/{id}
	// PUT localhost:8080/api/product/{id}
	// DELETE localhost:8080/api/product/{id}
	http.HandleFunc("/api/product/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProductByID(w,r)
		} else if r.Method == "PUT" {
			updateProductByID(w,r)
		} else if r.Method == "DELETE" {
			deleteProductByID(w,r)
		}
	})

	// GET localhost:8080/api/products
	// POST localhost:8080/api/product
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products)
		} else if r.Method == "POST" {
			var newProduct Product
			err := json.NewDecoder(r.Body).Decode(&newProduct)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// masuk ke dalam var product
			newProduct.ID = len(products) + 1
			products = append(products, newProduct)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newProduct)
		}
		
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"Message": "API Running",
		})
		// w.Write([]byte("OK"))
	}) 
	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}