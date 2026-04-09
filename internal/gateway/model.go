package gateway

import "net/http"

type ServiceName string

const (
	ServiceAuth ServiceName = "auth"
	ServiceUser ServiceName = "user"
	// ServiceProduct      ServiceName = "product"
	// ServiceOrder        ServiceName = "order"
	// ServicePayment      ServiceName = "payment"
	// ServiceNotification ServiceName = "notification"
)

type ServiceConfig struct {
	Name    ServiceName
	BaseURL string
}

type contextKey string

const (
	ContextKeyUserID    contextKey = "userId"
	ContextKeyRequestID contextKey = "requestID"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Service string `json:"service,omitempty"`
}

type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services,omitempty"`
}

type ProxyRequest struct {
	TargetService ServiceName
	StripPrefix   string
	OnError       func(w http.ResponseWriter, msg string)
}

type ProxyError struct {
	Code    string
	Message string
	Service ServiceName
}
