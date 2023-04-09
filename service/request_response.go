package service

type Field struct {
	Name     string   `json:"name"`
	Messages []string `json:"messages,omitempty"`
}

type Error struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Fields  []Field `json:"fields,omitempty"`
}

type BaseResponse struct {
	Error *Error `json:"error,omitempty"`
}

type OrderCreateRequest struct {
	MerchantId      string  `json:"merchant_id"`
	MerchantOrderId string  `json:"merchant_order_id"`
	UserId          string  `json:"user_id"`
	Amount          float32 `json:"amount"`
	Currency        string  `json:"currency"`
}

type OrderDetailResponse struct {
	BaseResponse
	Id              int64   `json:"id,omitempty"`
	TrackingId      string  `json:"tracking_id,omitempty"`
	MerchantId      string  `json:"merchant_id,omitempty"`
	MerchantOrderId string  `json:"merchant_order_id,omitempty"`
	UserId          string  `json:"user_id,omitempty"`
	Amount          float32 `json:"amount,omitempty"`
	Currency        string  `json:"currency,omitempty"`
}
