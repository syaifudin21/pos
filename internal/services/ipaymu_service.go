package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type IpaymuService struct {
	BaseURL string
	Va      string
	ApiKey  string
	DB      *gorm.DB
}

func NewIpaymuService(db *gorm.DB) *IpaymuService {
	return &IpaymuService{
		BaseURL: os.Getenv("IPAYMU_BASE_URL"),
		Va:      os.Getenv("IPAYMU_VA"),
		ApiKey:  os.Getenv("IPAYMU_API_KEY"),
		DB:      db,
	}
}

func (s *IpaymuService) header(body interface{}, method string) map[string]string {
	if method == "" {
		method = "POST"
	}
	var jsonBody string
	if body != nil {
		b, _ := json.Marshal(body)
		jsonBody = string(b)
	} else {
		jsonBody = "{}"
	}
	requestBody := strings.ToLower(fmt.Sprintf("%x", sha256.Sum256([]byte(jsonBody))))
	stringToSign := strings.ToUpper(method) + ":" + s.Va + ":" + requestBody + ":" + s.ApiKey
	h := hmac.New(sha256.New, []byte(s.ApiKey))
	h.Write([]byte(stringToSign))
	signature := hex.EncodeToString(h.Sum(nil))
	timestamp := time.Now().Format("20060102150405")

	fmt.Println("Debug Signature:", signature) // Debugging output
	fmt.Println("Debug Timestamp:", timestamp) // Debugging output
	fmt.Println("stringToSign:", stringToSign) // Debugging output
	// Debugging output for signature and timestamp
	if body != nil {
		fmt.Printf("Debug Request Body: %s\n", jsonBody) // Debugging output
	}
	return map[string]string{
		"signature": signature,
		"timestamp": timestamp,
		"va":        s.Va,
	}
}

func (s *IpaymuService) send(endPoint string, body interface{}, contentType, method string) (map[string]interface{}, error) {
	headers := s.header(body, method)
	var reqBody io.Reader
	if contentType == "multipart/form-data" {
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)
		if m, ok := body.(map[string]interface{}); ok {
			for k, v := range m {
				switch val := v.(type) {
				case string:
					_ = writer.WriteField(k, val)
				case *os.File:
					part, _ := writer.CreateFormFile(k, val.Name())
					io.Copy(part, val)
				}
			}
		}
		writer.Close()
		reqBody = buf
		contentType = writer.FormDataContentType()
	} else {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	}
	if method == "" {
		method = "POST"
	}
	url := s.BaseURL + endPoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("signature", headers["signature"])
	req.Header.Set("va", headers["va"])
	req.Header.Set("timestamp", headers["timestamp"])
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Tambahkan debug: baca body response mentah jika gagal decode JSON
	var res map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	decErr := json.Unmarshal(bodyBytes, &res)
	if decErr != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v, raw response: %s", decErr, string(bodyBytes))
	}
	return res, nil
}

// CreateDirectPayment creates a direct payment request to Ipaymu
func (s *IpaymuService) CreateDirectPayment(
	ServiceName string,
	ServiceRefID string,
	product []string,
	qty []int,
	price []int,
	name, email, phone, method, channel string,
	account *string,
) (map[string]interface{}, error) {
	// Hitung total amount dari seluruh produk
	amount := 0
	for i := range product {
		if i < len(qty) && i < len(price) {
			amount += qty[i] * price[i]
		}
	}
	body := map[string]interface{}{
		"product":        product,
		"qty":            qty,
		"price":          price,
		"name":           name,
		"email":          email,
		"phone":          phone,
		"expired":        24,
		"expiredType":    "hours",
		"referenceId":    1,
		"returnUrl":      os.Getenv("IPAYMU_RETURN_URL"), // set if needed
		"notifyUrl":      os.Getenv("IPAYMU_NOTIFY_URL"),
		"amount":         amount,
		"paymentMethod":  method,
		"paymentChannel": channel,
		"feeDirection":   "BUYER",
	}
	if account != nil {
		body["account"] = *account
	}
	endPoint := "/api/v2/payment/direct"

	res, err := s.send(endPoint, body, "application/json", "POST")
	if err != nil {
		return nil, err
	}

	// Ambil referenceIpaymu, totalStr, reqBodyStr, bodyBytes dari response
	var referenceIpaymu string
	var totalStr string

	// Ambil referenceIpaymu dan totalStr dari response
	if data, ok := res["Data"].(map[string]interface{}); ok {
		if ref, ok := data["TransactionId"]; ok {
			referenceIpaymu = fmt.Sprintf("%v", ref)
		}
		if total, ok := data["Total"]; ok {
			totalStr = fmt.Sprintf("%v", total)
		}
	}

	log := models.IpaymuLog{
		ServiceName:     ServiceName,
		ServiceRefID:    ServiceRefID,
		PaymentMethod:   method,
		PaymentChannel:  channel,
		RequestAt:       time.Now(),
		ReferenceIpaymu: referenceIpaymu,
	}
	// Ubah totalStr ke float64
	var amountFloat float64
	if totalStr != "" {
		fmt.Sscanf(totalStr, "%f", &amountFloat)
	}
	log.Amount = amountFloat
	if s.DB != nil {
		s.DB.Create(&log)
	}

	return res, nil

}

// UpdateIpaymuLogStatus updates the IpaymuLog status and timestamps based on notification
func (s *IpaymuService) NotifyDirectPayment(TrxId int, Status string, SettlementStatus string) error {
	// INSERT_YOUR_CODE
	var log models.IpaymuLog
	if err := s.DB.Where("reference_ipaymu = ?", fmt.Sprintf("%v", TrxId)).First(&log).Error; err != nil {
		return err
	}

	fmt.Println("Debug Signature:", log) // Debugging output

	dateNow := time.Now()

	updated := false

	if Status == "berhasil" {
		log.SuccessAt = &dateNow
		updated = true
	}
	if SettlementStatus == "settled" {
		log.SettlementAt = &dateNow
		updated = true
	}

	if updated {
		if err := s.DB.Save(&log).Error; err != nil {
			return err
		}
	}

	return nil
}

// Register melakukan pendaftaran user ke Ipaymu
func (s *IpaymuService) Register(
	name string,
	phone string,
	password string,
	email *string,
	optional map[string]interface{},
) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"name":     name,
		"phone":    phone,
		"password": password,
	}

	if email != nil {
		body["email"] = *email
		body["withoutEmail"] = "0"
	} else {
		body["withoutEmail"] = "1"
	}

	// Tambahkan field opsional jika ada
	if v, ok := optional["identityNo"]; ok {
		body["identityNo"] = v
	}
	if v, ok := optional["businessName"]; ok {
		body["businessName"] = v
	}
	if v, ok := optional["birthday"]; ok {
		body["birthday"] = v
	}
	if v, ok := optional["birthplace"]; ok {
		body["birthplace"] = v
	}
	if v, ok := optional["gender"]; ok {
		body["gender"] = v
	}
	if v, ok := optional["address"]; ok {
		body["address"] = v
	}

	var contentType string
	// Cek apakah ada identityPhoto (file)
	if v, ok := optional["identityPhoto"]; ok && v != nil {
		// v harus berupa *os.File
		body["identityPhoto"] = v
		contentType = "multipart/form-data"
	} else {
		contentType = "application/json"
	}

	res, err := s.send("/api/v2/register", body, contentType, "POST")
	if err != nil {
		return nil, err
	}

	if res != nil && err == nil {
		// Pastikan response mengandung data VA
		data, ok := res["Data"].(map[string]interface{})
		if ok {
			va, _ := data["Va"].(string)
			userIpaymu := &models.UserIpaymu{
				Name:      name,
				VaIpaymu:  va,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			// Simpan user_id jika ada di optional
			if v, ok := optional["user_id"]; ok {
				if userID, ok := v.(uint); ok {
					userIpaymu.UserID = userID
				}
			}
			// Simpan phone dan email jika ada
			if phone != "" {
				userIpaymu.Phone = &phone
			}
			if email != nil {
				userIpaymu.Email = email
			}
			// Simpan ke database
			s.DB.Create(userIpaymu)
		}
	}

	return res, nil
}
