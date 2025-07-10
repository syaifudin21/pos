package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type TsmService struct {
	DB                 *gorm.DB
	UserContextService *UserContextService
	UserPaymentService *UserPaymentService
}

func NewTsmService(db *gorm.DB, userContextService *UserContextService, userPaymentService *UserPaymentService) *TsmService {
	return &TsmService{
		DB:                 db,
		UserContextService: userContextService,
		UserPaymentService: userPaymentService,
	}
}

func (s *TsmService) RegisterTsm(userID uint, req dtos.TsmRegisterRequest) error {
	// If VaIpaymu is not provided in the request, try to get it from UserIpaymu
	if req.VaIpaymu == "" {
		va, err := s.UserPaymentService.GetUserIpaymuVa(userID)
		if err != nil {
			// If user has no iPaymu connection or other error, return error
			return errors.New("ipaymu_va_not_found_or_provided")
		}
		req.VaIpaymu = va
	}

	var existingTsm models.UserTsm
	result := s.DB.Where("user_id = ?", userID).First(&existingTsm)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new entry
			userTsm := models.UserTsm{
				UserID:       userID,
				AppCode:      req.AppCode,
				MerchantCode: req.MerchantCode,
				TerminalCode: req.TerminalCode,
				SerialNumber: req.SerialNumber,
				MID:          req.MID,
				VaIpaymu:     req.VaIpaymu,
			}
			if err := s.DB.Create(&userTsm).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}
	} else {
		// Update existing entry
		existingTsm.AppCode = req.AppCode
		existingTsm.MerchantCode = req.MerchantCode
		existingTsm.TerminalCode = req.TerminalCode
		existingTsm.SerialNumber = req.SerialNumber
		existingTsm.MID = req.MID
		existingTsm.VaIpaymu = req.VaIpaymu
		if err := s.DB.Save(&existingTsm).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *TsmService) base64urlEncode(data []byte) string {
	encoded := base64.RawURLEncoding.EncodeToString(data)
	return encoded
}

func (s *TsmService) generateHeader(bodyPart string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	jsonHeader, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	headerPart := s.base64urlEncode(jsonHeader)

	headerBodyPart := headerPart + "." + bodyPart
	key := os.Getenv("TSM_KEY")

	if key == "" {
		return "", errors.New("TSM_KEY environment variable not set")
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(headerBodyPart))
	barerPart := h.Sum(nil)

	signature := s.base64urlEncode(barerPart)
	headerToken := headerPart + "." + bodyPart + "." + signature

	return headerToken, nil
}

func (s *TsmService) GenerateAPPLink(req dtos.TsmGenerateApplinkRequest) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"app_code":       req.AppCode,
		"amount":         strconv.FormatFloat(req.Amount, 'f', -1, 64),
		"partner_trx_id": req.TrxID,
		"terminal_code":  req.TerminalCode,
		"merchant_code":  req.MerchantCode,
		"payment_method": "CARD",
		"timestamp":      time.Now().UnixNano() / int64(time.Millisecond),
	}

	jsonBodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	bodyPart := s.base64urlEncode(jsonBodyBytes)
	headerToken, err := s.generateHeader(bodyPart)
	if err != nil {
		return nil, err
	}

	baseUrl := os.Getenv("TSM_BASE")
	if baseUrl == "" {
		baseUrl = "https://tph-sandbox.tsmdev.id" // Default value from PHP code
	}

	client := &http.Client{}
	reqBody := bytes.NewBuffer(jsonBodyBytes)
	httpReq, err := http.NewRequest("POST", baseUrl+"/tph/v1/applink", reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Token", headerToken)

	// Log request to terminal
	for key, values := range httpReq.Header {
		for _, value := range values {
			log.Printf("  %s: %s", key, value)
		}
	}

	var reqBodyMap map[string]interface{}
	err = json.Unmarshal(jsonBodyBytes, &reqBodyMap)

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("TSM API Request Error: %v ", err)
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			log.Printf("  %s: %s", key, value)
		}
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("TSM API Response Body Read Error: %v ", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Pretty print response body
	var res map[string]interface{}
	json.Unmarshal(respBody, &res) // Unmarshal to pretty print, ignore error for logging

	if err := json.Unmarshal(respBody, &res); err != nil {
		log.Printf("TSM API Response Unmarshal Error: %v ", err)
		return nil, fmt.Errorf("failed to unmarshal response JSON: %w", err)
	}

	return res, nil
}
