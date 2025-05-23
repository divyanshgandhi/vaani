package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

// SMSProvider defines the interface for SMS providers
type SMSProvider interface {
	SendSMS(phoneNumber, message string) error
}

// SMSService handles SMS operations
type SMSService struct {
	provider SMSProvider
	logger   *zap.Logger
}

// NewSMSService creates a new SMS service with the appropriate provider
func NewSMSService(logger *zap.Logger) *SMSService {
	var provider SMSProvider

	// Check environment variables to determine which provider to use
	twilioSID := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioFromNumber := os.Getenv("TWILIO_FROM_NUMBER")

	if twilioSID != "" && twilioToken != "" && twilioFromNumber != "" {
		provider = &TwilioProvider{
			AccountSID: twilioSID,
			AuthToken:  twilioToken,
			FromNumber: twilioFromNumber,
			logger:     logger,
		}
		logger.Info("Using Twilio SMS provider")
	} else {
		provider = &MockSMSProvider{logger: logger}
		logger.Info("Using Mock SMS provider (for development)")
	}

	return &SMSService{
		provider: provider,
		logger:   logger,
	}
}

// SendOTP sends an OTP SMS to the given phone number
func (s *SMSService) SendOTP(phoneNumber, otp string) error {
	message := fmt.Sprintf("Your Vaani verification code is: %s. This code will expire in 5 minutes.", otp)

	s.logger.Info("Sending OTP SMS",
		zap.String("phone", phoneNumber),
		zap.String("otp", otp))

	return s.provider.SendSMS(phoneNumber, message)
}

// TwilioProvider implements SMS sending via Twilio
type TwilioProvider struct {
	AccountSID string
	AuthToken  string
	FromNumber string
	logger     *zap.Logger
}

// SendSMS sends an SMS using Twilio API
func (t *TwilioProvider) SendSMS(phoneNumber, message string) error {
	// Twilio API endpoint
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.AccountSID)

	// Prepare form data
	data := url.Values{}
	data.Set("From", t.FromNumber)
	data.Set("To", phoneNumber)
	data.Set("Body", message)

	// Create request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(t.AccountSID, t.AuthToken)

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return fmt.Errorf("twilio API error: %d - %v", resp.StatusCode, errorResp)
	}

	t.logger.Info("SMS sent successfully via Twilio",
		zap.String("phone", phoneNumber),
		zap.Int("status_code", resp.StatusCode))

	return nil
}

// MockSMSProvider implements a mock SMS provider for development
type MockSMSProvider struct {
	logger *zap.Logger
}

// SendSMS logs the SMS instead of actually sending it
func (m *MockSMSProvider) SendSMS(phoneNumber, message string) error {
	m.logger.Info("📱 MOCK SMS SENT",
		zap.String("phone", phoneNumber),
		zap.String("message", message))

	// In mock mode, also print to console for easy visibility
	fmt.Printf("\n🔔 SMS TO %s: %s\n\n", phoneNumber, message)

	return nil
}
