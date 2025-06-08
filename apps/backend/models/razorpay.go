package models

import "time"

type RazorpayOrderRequest struct {
	Amount          int               `json:"amount"`
	Currency        string            `json:"currency"`
	Receipt         string            `json:"receipt"`
	Notes           map[string]string `json:"notes,omitempty"`
	PartialPayment  bool              `json:"partial_payment,omitempty"`
}

type RazorpayOrderResponse struct {
	ID              string            `json:"id"`
	Entity          string            `json:"entity"`
	Amount          int               `json:"amount"`
	AmountPaid      int               `json:"amount_paid"`
	AmountDue       int               `json:"amount_due"`
	Currency        string            `json:"currency"`
	Receipt         string            `json:"receipt"`
	OfferId         string            `json:"offer_id"`
	Status          string            `json:"status"`
	Attempts        int               `json:"attempts"`
	Notes           map[string]string `json:"notes"`
	CreatedAt       int64             `json:"created_at"`
}


type RazorpayPaymentRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	OrderID  string `json:"order_id"`
	Method   string `json:"method"`
}

type RazorpayPaymentResponse struct {
	ID               string            `json:"id"`
	Entity           string            `json:"entity"`
	Amount           int               `json:"amount"`
	Currency         string            `json:"currency"`
	Status           string            `json:"status"`
	OrderID          string            `json:"order_id"`
	InvoiceID        string            `json:"invoice_id"`
	International    bool              `json:"international"`
	Method           string            `json:"method"`
	AmountRefunded   int               `json:"amount_refunded"`
	RefundStatus     string            `json:"refund_status"`
	Captured         bool              `json:"captured"`
	Description      string            `json:"description"`
	CardID           string            `json:"card_id"`
	Bank             string            `json:"bank"`
	Wallet           string            `json:"wallet"`
	VPA              string            `json:"vpa"`
	Email            string            `json:"email"`
	Contact          string            `json:"contact"`
	Notes            map[string]string `json:"notes"`
	Fee              int               `json:"fee"`
	Tax              int               `json:"tax"`
	ErrorCode        string            `json:"error_code"`
	ErrorDescription string            `json:"error_description"`
	CreatedAt        int64             `json:"created_at"`
}


type RazorpayWebhook struct {
	Entity    string                 `json:"entity"`
	AccountID string                 `json:"account_id"`
	Event     string                 `json:"event"`
	Contains  []string               `json:"contains"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt int64                  `json:"created_at"`
}

type RazorpayError struct {
	Error struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		Source      string `json:"source"`
		Step        string `json:"step"`
		Reason      string `json:"reason"`
		Metadata    struct {
			PaymentID string `json:"payment_id"`
			OrderID   string `json:"order_id"`
		} `json:"metadata"`
	} `json:"error"`
}


type PaymentVerificationRequest struct {
	OrderID   string `json:"razorpay_order_id" binding:"required"`
	PaymentID string `json:"razorpay_payment_id" binding:"required"`
	Signature string `json:"razorpay_signature" binding:"required"`
}

type CreateOrderRequest struct {
	Amount      float64           `json:"amount" binding:"required,gt=0"`
	Currency    string            `json:"currency"`
	Receipt     string            `json:"receipt"`
	Notes       map[string]string `json:"notes"`
	Description string            `json:"description"`
}


type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Details string      `json:"details,omitempty"`
}


type RazorpayRefundRequest struct {
	Amount int               `json:"amount,omitempty"`
	Speed  string            `json:"speed,omitempty"`
	Notes  map[string]string `json:"notes,omitempty"`
	Receipt string           `json:"receipt,omitempty"`
}

type RazorpayRefundResponse struct {
	ID        string            `json:"id"`
	Entity    string            `json:"entity"`
	Amount    int               `json:"amount"`
	Currency  string            `json:"currency"`
	PaymentID string            `json:"payment_id"`
	Notes     map[string]string `json:"notes"`
	Receipt   string            `json:"receipt"`
	Status    string            `json:"status"`
	CreatedAt int64             `json:"created_at"`
	BatchID   string            `json:"batch_id"`
}

const (
	OrderStatusCreated   = "created"
	OrderStatusAttempted = "attempted"
	OrderStatusPaid      = "paid"

	PaymentStatusCreated    = "created"
	PaymentStatusAuthorized = "authorized"
	PaymentStatusCaptured   = "captured"
	PaymentStatusRefunded   = "refunded"
	PaymentStatusFailed     = "failed"

	EventPaymentAuthorized = "payment.authorized"
	EventPaymentCaptured   = "payment.captured"
	EventPaymentFailed     = "payment.failed"
	EventOrderPaid         = "order.paid"
	EventRefundCreated     = "refund.created"
	EventRefundProcessed   = "refund.processed"
)

func (r *RazorpayOrderResponse) GetCreatedTime() time.Time {
	return time.Unix(r.CreatedAt, 0)
}

func (r *RazorpayPaymentResponse) GetCreatedTime() time.Time {
	return time.Unix(r.CreatedAt, 0)
}

func (r *RazorpayOrderResponse) IsCompleted() bool {
	return r.Status == OrderStatusPaid
}

func (r *RazorpayPaymentResponse) IsSuccessful() bool {
	return r.Status == PaymentStatusCaptured || r.Status == PaymentStatusAuthorized
}

func (r *RazorpayPaymentResponse) IsFailed() bool {
	return r.Status == PaymentStatusFailed
}

func ConvertToPaise(amount float64) int {
	return int(amount * 100)
}

func ConvertToRupees(amount int) float64 {
	return float64(amount) / 100
}