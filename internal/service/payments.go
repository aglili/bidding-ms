package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aglili/auction-app/internal/config"
	"github.com/aglili/auction-app/internal/domain"
)

type PaymentService struct {
	cfg *config.Config
}

func NewPaymentService(cfg *config.Config) *PaymentService {
	return &PaymentService{cfg: cfg}
}

const BASE_URL = "https://api.paystack.co"

func (s *PaymentService) InitializePayment(ctx context.Context, email string, amount int64, reference string) (*domain.PaymentResponse, error) {
	payload := domain.PaymentRequest{
		Email:     email,
		Amount:    amount,
		Reference: reference,
		Channels: []string{"card","bank_transfer","apple_pay","mobile_money","qr"},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", BASE_URL+"/transaction/initialize", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.cfg.PaystackSecretKey)
	req.Header.Set("Content-Type", "applicatication/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("paystack API error: %s - %s", resp.Status, string(body))
	}

	paymentResponse := &domain.PaymentResponse{}
	if err := json.Unmarshal(body, paymentResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !paymentResponse.Status {
		return nil, fmt.Errorf("payment initialization failed: %s", paymentResponse.Message)
	}

	return paymentResponse, nil
}
