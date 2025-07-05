package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
)

type IpaymuService struct {
	BaseURL    string
	APIKey     string
	VirtualAccount string
}

func NewIpaymuService() *IpaymuService {
	return &IpaymuService{
		BaseURL:    os.Getenv("IPAYMU_BASE_URL"),
		APIKey:     os.Getenv("IPAYMU_API_KEY"),
		VirtualAccount: os.Getenv("IPAYMU_VA"),
	}
}

func (s *IpaymuService) GenerateSignature(body string) string {
	stringToSign := fmt.Sprintf("%s:%s", s.VirtualAccount, body)
	hash := sha256.Sum256([]byte(stringToSign))
	return hex.EncodeToString(hash[:])
}

func (s *IpaymuService) CreateTransaction(req *dtos.IpaymuRequest) (*dtos.IpaymuResponse, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	signature := s.GenerateSignature(string(jsonBody))

	client := &http.Client{}
	request, err := http.NewRequest("POST", s.BaseURL+"/payment", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("VA", s.VirtualAccount)
	request.Header.Set("Signature", signature)
	request.Header.Set("Allow-Origin", "*")

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("iPaymu API error: Status %d, Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("iPaymu API returned non-OK status: %d - %s", resp.StatusCode, string(body))
	}

	var ipaymuResp dtos.IpaymuResponse
	if err := json.Unmarshal(body, &ipaymuResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if ipaymuResp.Status != 200 {
		return nil, errors.New(ipaymuResp.Message)
	}

	return &ipaymuResp, nil
}

func (s *IpaymuService) CheckTransactionStatus(transactionId string) (*dtos.IpaymuResponse, error) {
	// This is a simplified example. iPaymu's actual status check might require different parameters.
	// Refer to iPaymu API documentation for exact implementation.

	requestBody := map[string]string{
		"transactionId": transactionId,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	signature := s.GenerateSignature(string(jsonBody))

	client := &http.Client{}
	request, err := http.NewRequest("POST", s.BaseURL+"/payment/status", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("VA", s.VirtualAccount)
	request.Header.Set("Signature", signature)

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var ipaymuResp dtos.IpaymuResponse
	if err := json.Unmarshal(body, &ipaymuResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &ipaymuResp, nil
}

func (s *IpaymuService) ProcessIpaymuPayment(order *models.Order, user *models.User, orderItems []models.OrderItem) (*dtos.IpaymuResponse, error) {
	var products []string
	var qtys []int
	var prices []int

	for _, item := range orderItems {
		products = append(products, item.Product.Name)
		qtys = append(qtys, int(item.Quantity))
		prices = append(prices, int(item.Price))
	}

	// Construct the iPaymu request
	ipaymuReq := &dtos.IpaymuRequest{
		Product:    products,
		Qty:        qtys,
		Price:      prices,
		ReturnUrl:  os.Getenv("IPAYMU_RETURN_URL"), // Example return URL
		CancelUrl:  os.Getenv("IPAYMU_CANCEL_URL"), // Example cancel URL
		NotifyUrl:  os.Getenv("IPAYMU_NOTIFY_URL"), // Example notify URL
		ReferenceId: order.Uuid.String(),
		BuyerName:  user.Username,
		BuyerEmail: user.Username + "@example.com", // Assuming email is username@example.com
		BuyerPhone: "08123456789", // Example phone number
		Udf1:       order.Uuid.String(),
	}

	// Call iPaymu API
	ipaymuResp, err := s.CreateTransaction(ipaymuReq)
	if err != nil {
		return nil, fmt.Errorf("iPaymu transaction failed: %w", err)
	}

	return ipaymuResp, nil
}
