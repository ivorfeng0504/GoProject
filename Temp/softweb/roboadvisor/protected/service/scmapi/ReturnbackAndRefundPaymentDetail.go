package scmapi

type ReturnbackAndRefundPaymentDetail struct {
	PaymentTypeCode string
	PaymentType     string
	OnlineTypeCode  string
	OnlineType      string
	PaymentCode     string
	PaymentPrice    float32
	PaymentDate     string
	CanUseMoney     float32
	PaymentSource   string
}
