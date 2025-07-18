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
	"github.com/msyaifudin/pos/pkg/elasticsearch"
	"gorm.io/gorm"
)

type TsmService struct {
	DB                  *gorm.DB
	UserContextService  *UserContextService
	UserPaymentService  *UserPaymentService
	TsmLogService       *TsmLogService
	OrderPaymentService *OrderPaymentService
}

func NewTsmService(db *gorm.DB, userContextService *UserContextService, userPaymentService *UserPaymentService, tsmLogService *TsmLogService, orderPaymentService *OrderPaymentService) *TsmService {
	return &TsmService{
		DB:                  db,
		UserContextService:  userContextService,
		UserPaymentService:  userPaymentService,
		TsmLogService:       tsmLogService,
		OrderPaymentService: orderPaymentService,
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

func (s *TsmService) generateAPPLinkRequest(userID uint, appCode string, amount float64, trxID string, terminalCode string, merchantCode string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"app_code":       appCode,
		"amount":         strconv.FormatFloat(amount, 'f', -1, 64),
		"partner_trx_id": trxID,
		"terminal_code":  terminalCode,
		"merchant_code":  merchantCode,
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

	requestTime := time.Now()

	resp, err := client.Do(httpReq)

	var respBody []byte

	var finalRes map[string]interface{}
	var finalErr error
	var isSuccess bool

	if err != nil {
		finalErr = fmt.Errorf("failed to send HTTP request: %w", err)
		isSuccess = false
	} else {
		defer resp.Body.Close()
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			finalErr = fmt.Errorf("failed to read response body: %w", err)
			isSuccess = false
		} else {
			if err := json.Unmarshal(respBody, &finalRes); err != nil {
				finalErr = fmt.Errorf("failed to unmarshal response JSON: %w", err)
				isSuccess = false
			} else {
				isSuccess = true
			}
		}
	}

	// Log to database
	if s.TsmLogService != nil {
		logEntry := &models.TsmLog{
			UserID:       userID,
			ServiceName:  "TSM_APPLINK",
			ServiceRefID: trxID,
			Response:     string(respBody),
			IsPaid:       isSuccess, // Set IsPaid based on API call success
			RequestAt:    requestTime,
		}
		if logErr := s.TsmLogService.CreateTsmLog(logEntry); logErr != nil {
			log.Printf("Failed to save TSM log: %v", logErr)
		}
		// Log successful response to Elasticsearch
		logData := elasticsearch.APILog{
			Method:     "POST",
			Path:       baseUrl + "/tph/v1/applink",
			Status:     http.StatusOK,
			DurationMs: time.Since(requestTime).Milliseconds(),
			Extra: map[string]interface{}{
				"request_payload":  string(jsonBodyBytes),
				"response_payload": string(respBody),
				"service_name":     "TSM_APPLINK",
				"service_ref_id":   trxID,
			},
		}
		elasticsearch.LogAPI("tsm_curl_logs", logData)
	} else {
		// Log error response to Elasticsearch
		logData := elasticsearch.APILog{
			Method:     "POST",
			Path:       baseUrl + "/tph/v1/applink",
			Status:     0, // No HTTP status if request failed
			DurationMs: time.Since(requestTime).Milliseconds(),
			Error:      finalErr.Error(),
			Extra: map[string]interface{}{
				"request_payload": string(jsonBodyBytes),
				"service_name":    "TSM_APPLINK",
				"service_ref_id":  trxID,
			},
		}
		elasticsearch.LogAPI("tsm_curl_logs", logData)
	}

	return finalRes, finalErr
}

func (s *TsmService) GenerateAPPLink(userID uint, req dtos.TsmGenerateApplinkRequest) (map[string]interface{}, error) {
	return s.generateAPPLinkRequest(userID, req.AppCode, req.Amount, req.TrxID, req.TerminalCode, req.MerchantCode)
}

func (s *TsmService) HandleCallback(req dtos.TsmCallbackRequest) error {
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	jsonReq, err := json.Marshal(req)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to marshal TSM callback request: %w", err)
	}

	now := time.Now()

	// Find the existing TsmLog entry for the initial applink request
	var tsmLog models.TsmLog
	var currentUserID uint
	findErr := tx.Where("service_ref_id = ? AND service_name = ?", req.PartnerTrxID, "TSM_APPLINK").First(&tsmLog).Error

	if findErr != nil {
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			// If no existing log found, try to get UserID from OrderPayment
			var orderPayment models.OrderPayment
			if err := tx.Where("uuid = ?", req.PartnerTrxID).First(&orderPayment).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to find order payment for TSM callback: %w", err)
			}
			var order models.Order
			if err := tx.Where("id = ?", orderPayment.OrderID).First(&order).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to find order for TSM callback: %w", err)
			}
			currentUserID = order.UserID // Assign UserID from the associated Order

			newLogEntry := &models.TsmLog{
				UserID:       currentUserID,
				ServiceName:  "TSM_CALLBACK",
				ServiceRefID: req.PartnerTrxID,
				Callback:     string(jsonReq),
				IsPaid:       req.Status == "PAID",
				CallbackAt:   &now,
				RequestAt:    now, // Set request_at to now for new callback entries
			}
			if logErr := s.TsmLogService.CreateTsmLog(newLogEntry); logErr != nil {
				log.Printf("Failed to save new TSM callback log: %v", logErr)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to find TSM log for callback: %w", findErr)
		}
	} else {
		// Update the existing TsmLog entry
		tsmLog.Callback = string(jsonReq)
		tsmLog.IsPaid = req.Status == "PAID"
		tsmLog.CallbackAt = &now
		currentUserID = tsmLog.UserID // Get UserID from the existing log
		if logErr := tx.Save(&tsmLog).Error; logErr != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update TSM log for callback: %w", logErr)
		}
	}

	if req.Status == "PAID" {
		if err := s.OrderPaymentService.UpdateOrderPaymentAndStatus(tx, req.PartnerTrxID, req.Amount); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
