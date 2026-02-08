package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}
		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s. Available: %d, Requested: %d", productName, stock, item.Quantity)
		}
  
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetDailyReport() (*models.DailyReport, error) {
    report := &models.DailyReport{}

    queryTotal := `
        SELECT 
            COALESCE(SUM(total_amount), 0), 
            COUNT(id) 
        FROM transactions 
        WHERE DATE(created_at) = CURRENT_DATE`

    err := repo.db.QueryRow(queryTotal).Scan(&report.TotalRevenue, &report.TotalTransaksi)
    if err != nil {
        return nil, err
    }

    queryBestSeller := `
        SELECT 
            p.name, 
            SUM(td.quantity) as total_qty
        FROM transaction_details td
        JOIN transactions t ON td.transaction_id = t.id
        JOIN products p ON td.product_id = p.id
        WHERE DATE(t.created_at) = CURRENT_DATE
        GROUP BY p.name
        ORDER BY total_qty DESC
        LIMIT 1`

    err = repo.db.QueryRow(queryBestSeller).Scan(&report.ProdukTerlaris.Name, &report.ProdukTerlaris.Qty)
    
    if err != nil {
        if err == sql.ErrNoRows {
            report.ProdukTerlaris = models.ProductBestSeller{Name: "-", Qty: 0}
        } else {
            return nil, err
        }
    }

    return report, nil
}