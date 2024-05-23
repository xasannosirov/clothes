package postgresql

import (
	"context"
	"errors"
	"product-service/internal/entity"
	"time"

	"github.com/Masterminds/squirrel"
)

func (u *productRepo) CreateOrder(ctx context.Context, req *entity.Order) (*entity.Order, error) {
	data := map[string]any{
		"id":         req.Id,
		"product_id": req.ProductID,
		"user_id":    req.UserID,
		"status":     req.Status,
	}

	query, args, err := u.db.Sq.Builder.Insert(u.orderTable).SetMap(data).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *productRepo) GetOrder(ctx context.Context, params map[string]string) (*entity.Order, error) {
	var (
		order entity.Order
	)

	queryBuilder := u.ordersSelectQueryPrefix()
	for key, value := range params {
		if key == "user_id" {
			queryBuilder = queryBuilder.Where(squirrel.Eq{key: value})
		}
	}
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	if err = u.db.QueryRow(ctx, query, args...).Scan(
		&order.Id,
		&order.ProductID,
		&order.UserID,
		&order.Status,
	); err != nil {
		return nil, err
	}

	return &order, nil
}

func (u *productRepo) DeleteOrder(ctx context.Context, params map[string]string) error {
	caluses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.orderTable).
		SetMap(caluses).
		Where(u.db.Sq.Equal("id", params["order_id"])).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return err
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no sql rows")
	}

	return nil
}

func (u *productRepo) UserOrderHistory(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error) {
	return nil, nil
}
