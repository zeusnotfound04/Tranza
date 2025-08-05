package dto


type SendNotificationRequest struct {
	UserID   string            `json:"user_id" validate:"required"`
	Type     string            `json:"type" validate:"required,oneof=email sms push"`
	Template string            `json:"template" validate:"required"`
	Data     map[string]string `json:"data"`
	Priority string            `json:"priority" validate:"oneof=low normal high"`
}

// Notification Preferences Request
type NotificationPreferencesRequest struct {
	EmailEnabled         *bool `json:"email_enabled"`
	SMSEnabled           *bool `json:"sms_enabled"`
	PushEnabled          *bool `json:"push_enabled"`
	TransactionAlerts    *bool `json:"transaction_alerts"`
	AIPaymentAlerts      *bool `json:"ai_payment_alerts"`
	LowBalanceAlerts     *bool `json:"low_balance_alerts"`
	SecurityAlerts       *bool `json:"security_alerts"`
	MarketingEmails      *bool `json:"marketing_emails"`
	WeeklyReports        *bool `json:"weekly_reports"`
	MonthlyStatements    *bool `json:"monthly_statements"`
}

// Notification Preferences Response
type NotificationPreferencesResponse struct {
	EmailEnabled         bool   `json:"email_enabled"`
	SMSEnabled           bool   `json:"sms_enabled"`
	PushEnabled          bool   `json:"push_enabled"`
	TransactionAlerts    bool   `json:"transaction_alerts"`
	AIPaymentAlerts      bool   `json:"ai_payment_alerts"`
	LowBalanceAlerts     bool   `json:"low_balance_alerts"`
	SecurityAlerts       bool   `json:"security_alerts"`
	MarketingEmails      bool   `json:"marketing_emails"`
	WeeklyReports        bool   `json:"weekly_reports"`
	MonthlyStatements    bool   `json:"monthly_statements"`
	UpdatedAt            string `json:"updated_at"`
}