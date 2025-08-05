package dto

// Razorpay Webhook Payment Data
type RazorpayWebhookPayment struct {
	Entity PaymentEntity `json:"entity"`
}

type PaymentEntity struct {
	ID          string                 `json:"id"`
	Entity      string                 `json:"entity"`
	Amount      int64                  `json:"amount"`
	Currency    string                 `json:"currency"`
	Status      string                 `json:"status"`
	OrderID     string                 `json:"order_id"`
	Method      string                 `json:"method"`
	Description string                 `json:"description"`
	Email       string                 `json:"email"`
	Contact     string                 `json:"contact"`
	CreatedAt   int64                  `json:"created_at"`
	Notes       map[string]interface{} `json:"notes"`
}

// Webhook Event Request
type WebhookEventRequest struct {
	Account    string                 `json:"account"`
	Event      string                 `json:"event"`
	Contains   []string               `json:"contains"`
	Payload    map[string]interface{} `json:"payload"`
	CreatedAt  int64                  `json:"created_at"`
}