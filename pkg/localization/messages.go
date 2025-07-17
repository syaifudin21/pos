package localization

import (
	"encoding/json"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load messages from JSON files if you have them
	// bundle.LoadMessageFile("active.en.json")
	// bundle.LoadMessageFile("active.id.json")
}

var messages = map[string]map[string]string{
	"invalid_input": {
		"en": "Invalid input provided.",
		"id": "Masukan tidak valid.",
	},
	"failed_to_get_user_id": {
		"en": "Failed to get user ID from context.",
		"id": "Gagal mendapatkan ID pengguna dari konteks.",
	},
	"failed_to_get_owner_id": {
		"en": "Failed to get owner ID from context.",
		"id": "Gagal mendapatkan ID pemilik dari konteks.",
	},
	"payment_method_activated_successfully": {
		"en": "Payment method activated successfully.",
		"id": "Metode pembayaran berhasil diaktifkan.",
	},
	"payment_method_deactivated_successfully": {
		"en": "Payment method deactivated successfully.",
		"id": "Metode pembayaran berhasil dinonaktifkan.",
	},
	"user_payments_listed_successfully": {
		"en": "User payment methods listed successfully.",
		"id": "Metode pembayaran pengguna berhasil didaftarkan.",
	},
	"ipaymu_registration_successful": {
		"en": "iPaymu registration successful. You can now activate iPaymu payment methods.",
		"id": "Pendaftaran iPaymu berhasil. Anda sekarang dapat mengaktifkan metode pembayaran iPaymu.",
	},
	"tsm_registered_successfully": {
		"en": "TSM registered successfully.",
		"id": "TSM berhasil didaftarkan.",
	},
	"validation_generic_field_failed": {
		"en": "Field '%s' failed on the '%s' tag.",
		"id": "Bidang '%s' gagal pada tag '%s'.",
	},
	"validation_generic_error": {
		"en": "Validation failed.",
		"id": "Validasi gagal.",
	},
	// TSM Validation Messages
	"validation_appcode_required": {
		"en": "App Code is required.",
		"id": "Kode Aplikasi wajib diisi.",
	},
	"validation_merchantcode_required": {
		"en": "Merchant Code is required.",
		"id": "Kode Merchant wajib diisi.",
	},
	"validation_terminalcode_required": {
		"en": "Terminal Code is required.",
		"id": "Kode Terminal wajib diisi.",
	},
	"validation_serialnumber_required": {
		"en": "Serial Number is required.",
		"id": "Nomor Seri wajib diisi.",
	},
	"validation_mid_required": {
		"en": "MID is required.",
		"id": "MID wajib diisi.",
	},

	// Common Validation Messages
	"name_required": {
		"en": "Name is required.",
		"id": "Nama wajib diisi.",
	},
	"email_required": {
		"en": "Email is required.",
		"id": "Email wajib diisi.",
	},
	"password_required": {
		"en": "Password is required.",
		"id": "Kata sandi wajib diisi.",
	},
	"Role_required": {
		"en": "Role is required.",
		"id": "Peran wajib diisi.",
	},
	"OutletID_required": {
		"en": "Outlet ID is required.",
		"id": "ID Outlet wajib diisi.",
	},
	"OldPassword_required": {
		"en": "Old password is required.",
		"id": "Kata sandi lama wajib diisi.",
	},
	"NewPassword_required": {
		"en": "New password is required.",
		"id": "Kata sandi baru wajib diisi.",
	},
	"OTP_required": {
		"en": "OTP is required.",
		"id": "OTP wajib diisi.",
	},
	"email_invalid": {
		"en": "Invalid email format.",
		"id": "Format email tidak valid.",
	},
	"password_strength": {
		"en": "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one digit, and one special character.",
		"id": "Kata sandi harus minimal 8 karakter, mengandung setidaknya satu huruf besar, satu huruf kecil, satu angka, dan satu karakter khusus.",
	},
	"OTP_invalid": {
		"en": "Invalid OTP.",
		"id": "OTP tidak valid.",
	},
	"NewEmail_invalid": {
		"en": "Invalid new email format.",
		"id": "Format email baru tidak valid.",
	},
	"Product_required": {
		"en": "Product is required.",
		"id": "Produk wajib diisi.",
	},
	"Qty_required": {
		"en": "Quantity is required.",
		"id": "Kuantitas wajib diisi.",
	},
	"Price_required": {
		"en": "Price is required.",
		"id": "Harga wajib diisi.",
	},
	"Callback_required": {
		"en": "Callback URL is required.",
		"id": "URL Callback wajib diisi.",
	},
	"Method_required": {
		"en": "Payment method is required.",
		"id": "Metode pembayaran wajib diisi.",
	},
	"Channel_required": {
		"en": "Payment channel is required.",
		"id": "Channel pembayaran wajib diisi.",
	},
	"ReferenceId_required": {
		"en": "Reference ID is required.",
		"id": "ID Referensi wajib diisi.",
	},
	"without_email_required": {
		"en": "Without email field is required.",
		"id": "Bidang tanpa email wajib diisi.",
	},
	"url_invalid": {
		"en": "Invalid URL format.",
		"id": "Format URL tidak valid.",
	},
	"validation_greater_than_zero": {
		"en": "Field must be greater than zero.",
		"id": "Bidang harus lebih besar dari nol.",
	},
	"validation_dive_required": {
		"en": "All items in the list are required.",
		"id": "Semua item dalam daftar wajib diisi.",
	},
	"validation_required_if": {
		"en": "Field is required based on condition.",
		"id": "Bidang wajib diisi berdasarkan kondisi.",
	},
	"validation_required_with": {
		"en": "Field is required when other fields are present.",
		"id": "Bidang wajib diisi jika bidang lain ada.",
	},
	"product_min_one": {
		"en": "At least one product is required.",
		"id": "Setidaknya satu produk wajib diisi.",
	},
	"qty_min_one": {
		"en": "Quantity must be at least one.",
		"id": "Kuantitas harus setidaknya satu.",
	},
	"price_min_one": {
		"en": "Price must be at least one.",
		"id": "Harga harus setidaknya satu.",
	},
	"validation_oneof_invalid": {
		"en": "Invalid value for field.",
		"id": "Nilai tidak valid untuk bidang.",
	},
	"outlet_uuid_required": {
		"en": "Outlet UUID is required.",
		"id": "UUID Outlet wajib diisi.",
	},
	"order_items_required": {
		"en": "Order items are required.",
		"id": "Item pesanan wajib diisi.",
	},
	"payment_method_required": {
		"en": "Payment method is required.",
		"id": "Metode pembayaran wajib diisi.",
	},
	"product_uuid_required": {
		"en": "Product UUID is required.",
		"id": "UUID Produk wajib diisi.",
	},
	"quantity_required": {
		"en": "Quantity is required.",
		"id": "Kuantitas wajib diisi.",
	},
	"product_name_required": {
		"en": "Product name is required.",
		"id": "Nama produk wajib diisi.",
	},
	"product_description_required": {
		"en": "Product description is required.",
		"id": "Deskripsi produk wajib diisi.",
	},
	"product_price_required": {
		"en": "Product price is required.",
		"id": "Harga produk wajib diisi.",
	},
	"product_sku_required": {
		"en": "Product SKU is required.",
		"id": "SKU produk wajib diisi.",
	},
	"product_type_required": {
		"en": "Product type is required.",
		"id": "Tipe produk wajib diisi.",
	},
	"supplier_uuid_required": {
		"en": "Supplier UUID is required.",
		"id": "UUID Pemasok wajib diisi.",
	},
	"purchase_items_required": {
		"en": "Purchase items are required.",
		"id": "Item pembelian wajib diisi.",
	},
	"main_product_uuid_required": {
		"en": "Main product UUID is required.",
		"id": "UUID Produk Utama wajib diisi.",
	},
	"component_uuid_required": {
		"en": "Component UUID is required.",
		"id": "UUID Komponen wajib diisi.",
	},
	"contact_required": {
		"en": "Contact is required.",
		"id": "Kontak wajib diisi.",
	},
	"OrderUuid_required": {
		"en": "Order UUID is required.",
		"id": "UUID Pesanan wajib diisi.",
	},
	"PaymentMethodID_required": {
		"en": "Payment Method ID is required.",
		"id": "ID Metode Pembayaran wajib diisi.",
	},
	"AmountPaid_required": {
		"en": "Amount Paid is required.",
		"id": "Jumlah Pembayaran wajib diisi.",
	},
	"CustomerName_required": {
		"en": "Customer Name is required.",
		"id": "Nama Pelanggan wajib diisi.",
	},
	"CustomerEmail_required": {
		"en": "Customer Email is required.",
		"id": "Email Pelanggan wajib diisi.",
	},
	"CustomerPhone_required": {
		"en": "Customer Phone is required.",
		"id": "Telepon Pelanggan wajib diisi.",
	},
	"ipaymu_va_already_registered": {
		"en": "Your iPaymu VA is already registered.",
		"id": "VA iPaymu Anda sudah terdaftar.",
	},
}

func GetLocalizedMessage(messageKey, lang string) string {
	log.Printf("Localization: Requesting messageKey='%s' for lang='%s'", messageKey, lang)

	// First, try to get the message from the hardcoded map
	if langMessages, ok := messages[messageKey]; ok {
		if msg, ok := langMessages[lang]; ok {
			log.Printf("Localization: Found message in map: '%s'", msg)
			return msg
		}
		// Fallback to English if specific language not found in map
		if msg, ok := langMessages["en"]; ok {
			log.Printf("Localization: Found English fallback in map: '%s'", msg)
			return msg
		}
	}

	// If not found in map, try with go-i18n bundle (if message files are loaded)
	localizer := i18n.NewLocalizer(bundle, lang)
	translatedMessage, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageKey,
	})
	if err == nil && translatedMessage != messageKey { // Check if translation actually happened
		log.Printf("Localization: Found message in bundle: '%s'", translatedMessage)
		return translatedMessage
	}

	log.Printf("Localization: MessageKey '%s' not found. Returning original key.", messageKey)
	return messageKey // Fallback to original messageKey if no translation found
}
