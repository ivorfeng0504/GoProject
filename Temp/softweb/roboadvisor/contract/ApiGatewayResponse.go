package contract

import "time"

// ApiGatewayResponse API网关响应对象
type ApiGatewayResponse struct {
	RetCode          int
	RetMsg           string
	LastLoadApisTime time.Time
	IntervalTime     int
	ContentType      string
	Message          string
}
