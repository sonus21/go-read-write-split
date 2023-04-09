package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/sonus21/db-read-write/pkg/database"
	"time"
)

func doesOrderExist(ctx context.Context, db *sql.DB, req *OrderCreateRequest) (bool, error) {
	cntResult := db.QueryRowContext(ctx, "SELECT COUNT(id) from orders WHERE merchant_id=? AND merchant_order_id=?", req.MerchantId, req.MerchantOrderId)
	err := cntResult.Err()
	if err != nil {
		return false, err
	}
	var count int
	err = cntResult.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func CreateOrder(ctx context.Context, req *OrderCreateRequest) (*OrderDetailResponse, error) {
	// get db from context
	db := database.FromContext(ctx)
	exist, err := doesOrderExist(ctx, db, req)
	if err != nil {
		return nil, err
	}
	if exist {
		return &OrderDetailResponse{BaseResponse: BaseResponse{Error: &Error{
			Code:    400,
			Message: "Duplicate Order",
			Fields:  nil,
		}}}, nil
	}
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
	id := int64(-1)
	if err != nil {
		_ = txn.Rollback()
	} else {
		id, _ = result.LastInsertId()
	}

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
	err := row.Err()
	if err != nil {
		return nil, err
	}
	resp := OrderDetailResponse{}
	err = row.Scan(&resp.Id, &resp.Amount, &resp.Currency, &resp.MerchantId, &resp.MerchantOrderId, &resp.TrackingId, &resp.UserId)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return &OrderDetailResponse{BaseResponse: BaseResponse{Error: &Error{Code: 404, Message: "order not found"}}}, nil
	}
	return &resp, err
}
