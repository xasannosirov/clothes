package postgresql

import (
	"context"
	"errors"
	"product-service/internal/entity"
	"time"

	"github.com/Masterminds/squirrel"
)

func (u *productRepo) SaveToBasket(ctx context.Context, req *entity.Basket) (*entity.Basket, error) {
	data := map[string]any{
		"id":         req.ID,
		"user_id":    req.UserID,
		"product_id": req.ProductID,
	}

	query, args, err := u.db.Sq.Builder.Insert(u.basketTable).SetMap(data).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *productRepo) UpdateBasket(ctx context.Context, basket *entity.Basket) (*entity.Basket, error) {
	data := map[string]any{
		"id":         basket.ID,
		"product_id": basket.ProductID,
		"user_id":    basket.UserID,
		"count":      basket.Count,
		"updated_at": time.Now().Format(time.RFC3339),
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.basketTable).
		SetMap(data).
		Where(squirrel.Eq{"id": basket.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	if commandTag.RowsAffected() == 0 {
		return nil, errors.New("no sql rows")
	}

	return basket, nil
}

func (u *productRepo) ListBaskets(ctx context.Context, req *entity.GetWithID) (*entity.ListBasket, error) {
	products := &entity.ListBasket{}

	queryBuilder := u.basketsSelectQueryPrefix()
	queryBuilder = queryBuilder.Where(squirrel.Eq{"id": req.ID})
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Basket
		if err = rows.Scan(
			&product.ID,
			&product.ProductID,
			&product.UserID,
		); err != nil {
			return nil, err
		}
		products.Baskets = append(products.Baskets, &product)

	}

	var count uint64
	total := "SELECT COUNT(*) FROM products WHERE deleted_at IS NULL"
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		products.TotalCount = 0
	}
	products.TotalCount = count

	return products, nil
}

func (u *productRepo) DeleteFromBasket(ctx context.Context, id string) error {
	caluses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.basketTable).
		SetMap(caluses).
		Where(u.db.Sq.Equal("id", id)).
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
		return err
	}

	return nil
}
