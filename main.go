package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Product struct {
	ID    int 		`json:"id"`
	Name  string 	`json:"name"`
	Price int 		`json:"price"`
	Stock int 		`json:"stock"`
}

type Category struct {
	ID    		int 		`json:"id"`
	Name  		string 	`json:"name"`
	Description string 	`json:"description"`
}

var products = []Product{
	{ID: 1, Name: "Product 1", Price: 1000, Stock: 10},
	{ID: 2, Name: "Product 2", Price: 2000, Stock: 20},
	{ID: 3, Name: "Product 3", Price: 3000, Stock: 30},
}

var categories = []Category{
	{ID: 1, Name: "Category 1", Description: "Description 1"},
	{ID: 2, Name: "Category 2", Description: "Description 2"},
	{ID: 3, Name: "Category 3", Description: "Description 3"},
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

// Function for category
func getAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
			Port:   viper.GetString("PORT"),
			DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer db.Close()
	
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo, categoryRepo)
	productHandler := handlers.NewProductHandler(productService)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

	http.HandleFunc("/api/categories", categoryHandler.HandleCategory)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)
	http.HandleFunc("/api/report/hari-ini", transactionHandler.HandleDailyReport)

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"Message": "API Running",
		})
		// w.Write([]byte("OK"))
	}) 
	fmt.Println("Server running di localhost:" + config.Port)

	err = http.ListenAndServe(":" + config.Port, nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}