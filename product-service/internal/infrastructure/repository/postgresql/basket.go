package postgresql

import (
	"context"
	"log"
	"product-service/internal/entity"

	"github.com/Masterminds/squirrel"
)

func (u *productRepo) SaveToBasket(ctx context.Context, req *entity.BasketCreateReq) (*entity.MoveResponse, error) {
	query := `SELECT COUNT(*) FROM baskets WHERE user_id = $1 AND product_id = $2`
	var count int
	if err := u.db.QueryRow(ctx, query, req.UserID, req.ProductID).Scan(&count); err != nil {
		return nil, err
	}

	var status bool
	if count == 0 {
		insertQuery := `INSERT INTO baskets (user_id, product_id) VALUES ($1, $2)`
		result, err := u.db.Exec(ctx, insertQuery, req.UserID, req.ProductID)
		if err != nil {
			return nil, err
		}

		if result.RowsAffected() == 0 {
			status = false
		} else {
			status = true
		}
	} else {
		deleteQuery := `DELETE FROM baskets WHERE user_id = $1 AND product_id = $2`
		result, err := u.db.Exec(ctx, deleteQuery, req.UserID, req.ProductID)
		if err != nil {
			return nil, err
		}

		if result.RowsAffected() == 0 {
			status = true
		} else {
			status = false
		}
	}

	return &entity.MoveResponse{
		Status: status,
	}, nil
}

func (u *productRepo) GetUserBaskets(ctx context.Context, req *entity.GetWithID) (*entity.ListProduct, error) {
	queryBuilder := u.db.Sq.Builder.Select("product_id").From(u.basketTable)
	queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": req.ID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response entity.ListProduct

	for rows.Next() {
		var (
			productID string
		)
		err := rows.Scan(&productID)
		if err != nil {
			return nil, err
		}

		product, err := u.GetProduct(ctx, map[string]string{
			"id": productID,
		})
		if err != nil {
			log.Println(err.Error(), "no append product to array in user baskets storage in product service")
		} else {
			response.Products = append(response.Products, product)
		}
	}

	return &response, nil
}
