package postgresql

import (
	"context"
	"database/sql"
	"product-service/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

func (u *productRepo) SaveToBasket(ctx context.Context, req *entity.BasketCreateReq) (*entity.Basket, error) {
	var existingProductIDs []string
	err := u.db.QueryRow(ctx, `SELECT product_id FROM `+u.basketTable+` WHERE user_id = $1`, req.UserID).Scan(pq.Array(&existingProductIDs))
	if err != nil && err != sql.ErrNoRows {
		return nil, err

	}

	if err == sql.ErrNoRows {
		productIDs := []string{req.ProductID}
		data := map[string]any{
			"user_id":    req.UserID,
			"product_id": pq.Array(productIDs),
		}

		query, args, err := u.db.Sq.Builder.Insert(u.basketTable).SetMap(data).ToSql()
		if err != nil {
			return nil, err
		}

		_, err = u.db.Exec(ctx, query, args...)
		if err != nil {
			return nil, err
		}
	} else {		
		var exists bool
			for _, existingProductID := range existingProductIDs {
				if existingProductID == req.ProductID {
					exists = true
					break
				}
			}

			if exists {
				updatedProductIDs := make([]string, 0, len(existingProductIDs))
				for _, pid := range existingProductIDs {
					if pid != req.ProductID {
						updatedProductIDs = append(updatedProductIDs, pid)
					}
				}

				data := map[string]any{
					"product_id": pq.Array(updatedProductIDs),
				}

				query, args, err := u.db.Sq.Builder.
					Update(u.basketTable).
					SetMap(data).
					Where("user_id = ?", req.UserID).
					ToSql()
				if err != nil {
					return nil, err
				}

				_, err = u.db.Exec(ctx, query, args...)
				if err != nil {
					return nil, err
				}

			} else {
				
				_, err = u.db.Exec(ctx, `UPDATE `+u.basketTable+` SET product_id = array_append(product_id, $2::uuid) WHERE user_id = $1`, req.UserID, req.ProductID)
				if err != nil {
					return nil, err
				}
			
		}
	}

	return &entity.Basket{
		UserID: req.UserID,
	}, nil
}

func (u *productRepo) GetBasket(ctx context.Context, req *entity.GetBAsketReq) (*entity.Basket, error) {
	product := &entity.Basket{}

	offset := (req.Page - 1 ) * req.Limit
	queryBuilder := u.basketsSelectQueryPrefix()
	queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": req.UserId})
	queryBuilder = queryBuilder.Limit(uint64(req.Limit))
	queryBuilder = queryBuilder.Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = u.db.QueryRow(ctx, query, args...).Scan(
		&product.ProductIDs,
		&product.UserID,
	)
	if err != nil {
		return nil, err
	}
	var existingProductIDs []string
	err = u.db.QueryRow(ctx, `SELECT product_id FROM `+u.basketTable+` WHERE user_id = $1`, req.UserId).Scan(pq.Array(&existingProductIDs))
	if err != nil {
		return nil, err
	}
	product.TotalCount = int64(len(existingProductIDs))
	return product, nil
}

func (u *productRepo) DeleteFromBasket(ctx context.Context, userID string, productID string) error {
	var existingProductIDs []string
	err := u.db.QueryRow(ctx, `SELECT product_id FROM `+u.basketTable+` WHERE user_id = $1`, userID).Scan(pq.Array(&existingProductIDs))
	if err != nil {
		return err
	}

	updatedProductIDs := make([]string, 0, len(existingProductIDs))
	for _, pid := range existingProductIDs {
		if pid != productID {
			updatedProductIDs = append(updatedProductIDs, pid)
		}
	}

	data := map[string]any{
		"product_id": pq.Array(updatedProductIDs),
	}

	query, args, err := u.db.Sq.Builder.
		Update(u.basketTable).
		SetMap(data).
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
