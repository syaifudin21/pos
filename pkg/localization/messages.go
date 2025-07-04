package localization

var messages = map[string]map[string]string{
	"invalid_request_payload": {
		"en": "Invalid request payload",
		"id": "Payload permintaan tidak valid",
	},
	"username_password_required": {
		"en": "Username and password are required",
		"id": "Nama pengguna dan kata sandi wajib diisi",
	},
	"invalid_role_specified": {
		"en": "Invalid role specified",
		"id": "Peran yang ditentukan tidak valid",
	},
	"user_registered_successfully": {
		"en": "User registered successfully",
		"id": "Pengguna berhasil didaftarkan",
	},
	"login_successful": {
		"en": "Login successful",
		"id": "Login berhasil",
	},
	"invalid_user_uuid_format": {
		"en": "Invalid User Uuid format",
		"id": "Format UUID Pengguna tidak valid",
	},
	"user_blocked_successfully": {
		"en": "User blocked successfully",
		"id": "Pengguna berhasil diblokir",
	},
	"user_unblocked_successfully": {
		"en": "User unblocked successfully",
		"id": "Pengguna berhasil dibuka blokirnya",
	},
	"outlet_not_found": {
		"en": "outlet not found",
		"id": "Outlet tidak ditemukan",
	},
	"product_not_found": {
		"en": "product not found",
		"id": "Produk tidak ditemukan",
	},
	"supplier_not_found": {
		"en": "supplier not found",
		"id": "Pemasok tidak ditemukan",
	},
	"recipe_not_found": {
		"en": "recipe not found",
		"id": "Resep tidak ditemukan",
	},
	"stock_not_found": {
		"en": "stock not found",
		"id": "Stok tidak ditemukan",
	},
	"order_not_found": {
		"en": "order not found",
		"id": "Pesanan tidak ditemukan",
	},
	"purchase_order_not_found": {
		"en": "purchase order not found",
		"id": "Pesanan pembelian tidak ditemukan",
	},
	"invalid_uuid_format": {
		"en": "Invalid UUID format",
		"id": "Format UUID tidak valid",
	},
	"order_created_successfully": {
		"en": "Order created successfully",
		"id": "Pesanan berhasil dibuat",
	},
	"order_retrieved_successfully": {
		"en": "Order retrieved successfully",
		"id": "Pesanan berhasil diambil",
	},
	"orders_retrieved_successfully": {
		"en": "Orders retrieved successfully",
		"id": "Pesanan berhasil diambil",
	},
	"outlets_retrieved_successfully": {
		"en": "Outlets retrieved successfully",
		"id": "Outlet berhasil diambil",
	},
	"outlet_created_successfully": {
		"en": "Outlet created successfully",
		"id": "Outlet berhasil dibuat",
	},
	"outlet_updated_successfully": {
		"en": "Outlet updated successfully",
		"id": "Outlet berhasil diperbarui",
	},
	"outlet_deleted_successfully": {
		"en": "Outlet deleted successfully",
		"id": "Outlet berhasil dihapus",
	},
	"products_retrieved_successfully": {
		"en": "Products retrieved successfully",
		"id": "Produk berhasil diambil",
	},
	"product_created_successfully": {
		"en": "Product created successfully",
		"id": "Produk berhasil dibuat",
	},
	"product_updated_successfully": {
		"en": "Product updated successfully",
		"id": "Produk berhasil diperbarui",
	},
	"product_deleted_successfully": {
		"en": "Product deleted successfully",
		"id": "Produk berhasil dihapus",
	},
	"invalid_product_type_specified": {
		"en": "Invalid product type specified",
		"id": "Jenis produk yang ditentukan tidak valid",
	},
	"purchase_order_created_successfully": {
		"en": "Purchase order created successfully",
		"id": "Pesanan pembelian berhasil dibuat",
	},
	"purchase_order_retrieved_successfully": {
		"en": "Purchase order retrieved successfully",
		"id": "Pesanan pembelian berhasil diambil",
	},
	"purchase_orders_retrieved_successfully": {
		"en": "Purchase orders retrieved successfully",
		"id": "Pesanan pembelian berhasil diambil",
	},
	"purchase_order_received_successfully": {
		"en": "Purchase order received successfully",
		"id": "Pesanan pembelian berhasil diterima",
	},
	"invalid_main_product_uuid_format": {
		"en": "Invalid Main Product Uuid format",
		"id": "Format UUID Produk Utama tidak valid",
	},
	"recipes_retrieved_successfully": {
		"en": "Recipes retrieved successfully",
		"id": "Resep berhasil diambil",
	},
	"recipe_created_successfully": {
		"en": "Recipe created successfully",
		"id": "Resep berhasil dibuat",
	},
	"recipe_updated_successfully": {
		"en": "Recipe updated successfully",
		"id": "Resep berhasil diperbarui",
	},
	"recipe_deleted_successfully": {
		"en": "Recipe deleted successfully",
		"id": "Resep berhasil dihapus",
	},
	"stock_retrieved_successfully": {
		"en": "Stock retrieved successfully",
		"id": "Stok berhasil diambil",
	},
	"outlet_stocks_retrieved_successfully": {
		"en": "Outlet stocks retrieved successfully",
		"id": "Stok outlet berhasil diambil",
	},
	"stock_updated_successfully": {
		"en": "Stock updated successfully",
		"id": "Stok berhasil diperbarui",
	},
	"sales_report_by_outlet_generated_successfully": {
		"en": "Sales report by outlet generated successfully",
		"id": "Laporan penjualan berdasarkan outlet berhasil dibuat",
	},
	"sales_report_by_product_generated_successfully": {
		"en": "Sales report by product generated successfully",
		"id": "Laporan penjualan berdasarkan produk berhasil dibuat",
	},
	"invalid_start_date_format": {
		"en": "Invalid start_date format. Use YYYY-MM-DD",
		"id": "Format tanggal mulai tidak valid. Gunakan YYYY-MM-DD",
	},
	"invalid_end_date_format": {
		"en": "Invalid end_date format. Use YYYY-MM-DD",
		"id": "Format tanggal akhir tidak valid. Gunakan YYYY-MM-DD",
	},
	"invalid_outlet_uuid_format": {
		"en": "Invalid Outlet Uuid format",
		"id": "Format UUID Outlet tidak valid",
	},
	"invalid_product_uuid_format": {
		"en": "Invalid Product Uuid format",
		"id": "Format UUID Produk tidak valid",
	},
	"user_deleted_successfully": {
		"en": "User deleted successfully",
		"id": "Pengguna berhasil dihapus",
	},
}

// GetLocalizedMessage retrieves a message for a given key and language.
// It falls back to English if the requested language is not found.
func GetLocalizedMessage(key, lang string) string {
	if langMessages, ok := messages[key]; ok {
		if msg, ok := langMessages[lang]; ok {
			return msg
		}
		// Fallback to English if specific language not found for the key
		if msg, ok := langMessages["en"]; ok {
			return msg
		}
	}
	// Return the key itself if no translation is found
	return key
}

var validationMessages = map[string]map[string]string{
	"name_required": {
		"en": "Name is required",
		"id": "Nama wajib diisi",
	},
	"contact_required": {
		"en": "Contact is required",
		"id": "Kontak wajib diisi",
	},
	"address_required": {
		"en": "Address is required",
		"id": "Alamat wajib diisi",
	},
	"product_name_required": {
		"en": "Product name is required",
		"id": "Nama produk wajib diisi",
	},
	"product_description_required": {
		"en": "Product description is required",
		"id": "Deskripsi produk wajib diisi",
	},
	"product_price_required": {
		"en": "Product price is required",
		"id": "Harga produk wajib diisi",
	},
	"product_sku_required": {
		"en": "Product SKU is required",
		"id": "SKU produk wajib diisi",
	},
	"product_type_required": {
		"en": "Product type is required and must be one of retail_item, fnb_main_product, or fnb_component",
		"id": "Tipe produk wajib diisi dan harus salah satu dari retail_item, fnb_main_product, atau fnb_component",
	},
	"username_required": {
		"en": "Username is required",
		"id": "Username wajib diisi",
	},
	"password_required": {
		"en": "Password is required",
		"id": "Password wajib diisi",
	},
	"role_required": {
		"en": "Role is required",
		"id": "Role wajib diisi",
	},
	"outlet_id_required": {
		"en": "Outlet ID is required",
		"id": "ID Outlet wajib diisi",
	},
	"email_required": {
		"en": "Email is required",
		"id": "Email wajib diisi",
	},
	"phone_number_required": {
		"en": "Phone number is required",
		"id": "Nomor telepon wajib diisi",
	},
	"outlet_uuid_required": {
		"en": "Outlet UUID is required",
		"id": "UUID Outlet wajib diisi",
	},
	"order_items_required": {
		"en": "Order items are required",
		"id": "Item pesanan wajib diisi",
	},
	"payment_method_required": {
		"en": "Payment method is required",
		"id": "Metode pembayaran wajib diisi",
	},
	"product_uuid_required": {
		"en": "Product UUID is required",
		"id": "UUID Produk wajib diisi",
	},
	"quantity_required": {
		"en": "Quantity is required",
		"id": "Kuantitas wajib diisi",
	},
	"supplier_uuid_required": {
		"en": "Supplier UUID is required",
		"id": "UUID Supplier wajib diisi",
	},
	"purchase_items_required": {
		"en": "Purchase items are required",
		"id": "Item pembelian wajib diisi",
	},
	"price_required": {
		"en": "Price is required",
		"id": "Harga wajib diisi",
	},
	"main_product_uuid_required": {
		"en": "Main Product UUID is required",
		"id": "UUID Produk Utama wajib diisi",
	},
	"component_uuid_required": {
		"en": "Component UUID is required",
		"id": "UUID Komponen wajib diisi",
	},
	"product_required": {
		"en": "Product is required",
		"id": "Produk wajib diisi",
	},
	"qty_required": {
		"en": "Quantity is required",
		"id": "Kuantitas wajib diisi",
	},
	"return_url_required": {
		"en": "Return URL is required",
		"id": "Return URL wajib diisi",
	},
	"cancel_url_required": {
		"en": "Cancel URL is required",
		"id": "Cancel URL wajib diisi",
	},
	"notify_url_required": {
		"en": "Notify URL is required",
		"id": "Notify URL wajib diisi",
	},
	"reference_id_required": {
		"en": "Reference ID is required",
		"id": "Reference ID wajib diisi",
	},
	"buyer_name_required": {
		"en": "Buyer name is required",
		"id": "Nama pembeli wajib diisi",
	},
	"buyer_email_required": {
		"en": "Buyer email is required",
		"id": "Email pembeli wajib diisi",
	},
	"buyer_phone_required": {
		"en": "Buyer phone is required",
		"id": "Nomor telepon pembeli wajib diisi",
	},
	"udf1_required": {
		"en": "UDF1 is required",
		"id": "UDF1 wajib diisi",
	},
	"role_invalid": {
		"en": "Role must be one of admin, owner, manager, or cashier",
		"id": "Role harus salah satu dari admin, owner, manager, atau cashier",
	},
	"password_strength": {
		"en": "Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character",
		"id": "Kata sandi harus mengandung setidaknya satu huruf kapital, satu huruf kecil, satu angka, dan satu karakter spesial",
	},
}

// GetLocalizedValidationMessage retrieves a validation message for a given key and language.
// It falls back to English if the requested language is not found.
func GetLocalizedValidationMessage(key, lang string) string {
	if langMessages, ok := validationMessages[key]; ok {
		if msg, ok := langMessages[lang]; ok {
			return msg
		}
		// Fallback to English if specific language not found for the key
		if msg, ok := langMessages["en"]; ok {
			return msg
		}
	}
	// Return the key itself if no translation is found
	return key
}
