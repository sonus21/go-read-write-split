package service

type OrderCreateRequest struct {
	MerchantId      string  `json:"merchant_id"`
	MerchantOrderId string  `json:"merchant_order_id"`
	UserId          string  `json:"user_id"`
	Amount          float32 `json:"amount"`
	Currency        string  `json:"currency"`
}

type OrderDetailResponse struct {
	Id              int64   `json:"id"`
	TrackingId      string  `json:"tracking_id"`
	MerchantId      string  `json:"merchant_id"`
	MerchantOrderId string  `json:"merchant_order_id"`
	UserId          string  `json:"user_id"`
	Amount          float32 `json:"amount"`
	Currency        string  `json:"currency"`
}
