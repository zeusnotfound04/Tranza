package razorpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// WebhookEvent represents a Razorpay webhook event
type WebhookEvent struct {
	Account   string                 `json:"account"`
	Event     string                 `json:"event"`
	Contains  []string               `json:"contains"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt int64                  `json:"created_at"`
}

// VerifyWebhookSignature verifies webhook signature
func (c *Client) VerifyWebhookSignature(body []byte, signature, secret string) bool {
	expectedSignature := generateWebhookSignature(body, secret)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// ParseWebhookEvent parses webhook event from JSON
func (c *Client) ParseWebhookEvent(body []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func generateWebhookSignature(body []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

// pkg/razorpay/config.go
