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
)

type IpaymuService struct {
	BaseURL string
	Va      string
	ApiKey  string
}

func NewIpaymuService() *IpaymuService {
	return &IpaymuService{
		BaseURL: os.Getenv("IPAYMU_BASE_URL"),
		Va:      os.Getenv("IPAYMU_VA"),
		ApiKey:  os.Getenv("IPAYMU_API_KEY"),
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
	product []string,
	qty []int,
	price []int,
	name, email, phone, callback, method, channel string,
	account *string,
) (map[string]interface{}, error) {
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
		"returnUrl":      os.Getenv("IPAYMU_NOTIFY_URL"), // set if needed
		"notifyUrl":      callback,
		"amount":         price[0] * qty[0],
		"paymentMethod":  method,
		"paymentChannel": channel,
		"feeDirection":   "BUYER",
	}
	if account != nil {
		body["account"] = *account
	}

	return s.send("/api/v2/payment/direct", body, "application/json", "POST")
}
