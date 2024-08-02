package postgresql

import (
	"context"
	"fmt"
	"product-service/internal/entity"
)

func (u *productRepo) SaveToBasket(ctx context.Context, req *entity.BasketCreateReq) (*entity.MoveResponse, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM baskets WHERE user_id = '%s' AND product_id = '%s'`, req.UserID, req.ProductID)
	var count int
	if err := u.db.QueryRow(ctx, query).Scan(&count); err != nil {
		return nil, err
	}

	var status bool
	if count == 0 {
		insertQuery := fmt.Sprintf(`INSERT INTO baskets (user_id, product_id) VALUES ('%s', '%s')`, req.UserID, req.ProductID)
		result, err := u.db.Exec(ctx, insertQuery)
		if err != nil {
			return nil, err
		}

		if result.RowsAffected() == 0 {
			status = false
		} else {
			status = true
		}
	} else {
		deleteQuery := fmt.Sprintf(`DELETE FROM baskets WHERE user_id = '%s' AND product_id = '%s'`, req.UserID, req.ProductID)
		result, err := u.db.Exec(ctx, deleteQuery)
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
	query := `
	SELECT 
	    p.id,
	    p.name,
	    c.name,
	    p.description,
	    p.made_in,
	    p.count,
	    p.cost,
	    p.discount,
	    p.color,
	    p.size,
	    p.age_min,
	    p.age_max,
	    p.for_gender
	FROM baskets AS b
	INNER JOIN products AS p ON b.product_id = p.id
	INNER JOIN users u on u.id = b.user_id
	INNER JOIN category c on c.id = p.category_id
	WHERE b.user_id = $1
	`

	var response entity.ListProduct

	rows, err := u.db.Query(ctx, query, req.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Category,
			&product.Description,
			&product.MadeIn,
			&product.Count,
			&product.Cost,
			&product.Discount,
			&product.Color,
			&product.Size,
			&product.AgeMin,
			&product.AgeMax,
			&product.ForGender,
		)
		if err != nil {
			return nil, err
		}

		response.Products = append(response.Products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Println("sql basket", response)
	return &response, nil
}
