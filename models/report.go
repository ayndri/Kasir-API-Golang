package models

type ProductBestSeller struct {
	Name string `json:"nama"`
	Qty  int    `json:"qty_terjual"`
}

type DailyReport struct {
	TotalRevenue   int               `json:"total_revenue"`
	TotalTransaksi int               `json:"total_transaksi"`
	ProdukTerlaris ProductBestSeller `json:"produk_terlaris"`
}