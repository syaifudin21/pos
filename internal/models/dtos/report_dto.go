package dtos

type StockReportResponse struct {
	ProductName string  `json:"product_name"`
	ProductSku  string  `json:"product_sku,omitempty"`
	VariantName string  `json:"variant_name,omitempty"`
	VariantSku  string  `json:"variant_sku,omitempty"`
	Quantity    float64 `json:"quantity"`
}
