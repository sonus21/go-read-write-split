package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/sonus21/db-read-write/pkg/database"
	"time"
)

func CreateOrder(ctx context.Context, req *OrderCreateRequest) (*OrderDetailResponse, error) {
	db := database.FromContext(ctx)
	txn, err := db.BeginTx(ctx, nil)
	trackingId := uuid.New().String()

	result, err := txn.ExecContext(ctx, "INSERT INTO orders  (amount, created_at, currency, merchant_id, merchant_order_id, tracking_id,updated_at, user_id) value(?,?,?,?,?,?,?,?)",
		req.Amount,
		time.Now(),
		req.Currency,
		req.MerchantId,
		req.MerchantOrderId,
		trackingId,
		time.Now(),
		req.UserId)
	if err == nil {
		err = txn.Commit()
	}
	if err != nil {
		_ = txn.Rollback()
	}
	id, _ := result.LastInsertId()
	return &OrderDetailResponse{
		Id:              id,
		TrackingId:      trackingId,
		MerchantId:      req.MerchantId,
		MerchantOrderId: req.MerchantOrderId,
		UserId:          req.UserId,
		Amount:          req.Amount,
		Currency:        req.Currency,
	}, err
}

func OrderDetail(ctx context.Context, orderId string) (*OrderDetailResponse, error) {
	db := database.FromContext(ctx)
	// use db to query the database
	row := db.QueryRowContext(ctx, "SELECT id, amount, currency, merchant_id, merchant_order_id, tracking_id, user_id FROM orders where id=?", orderId)
	// check for record not found etc
	if row.Err() != nil {
		return nil, row.Err()
	}
	resp := OrderDetailResponse{}
	err := row.Scan(&resp.Id, &resp.Amount, &resp.Currency, &resp.MerchantId, &resp.MerchantOrderId, &resp.TrackingId, &resp.UserId)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
