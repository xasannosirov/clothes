package postgresql

import (
	"context"
	"errors"
	"product-service/internal/entity"

	"github.com/Masterminds/squirrel"
)

func (p *productRepo) IsUnique(ctx context.Context, tableName, UserId, ProductId string) (bool, error) {

	queryBuilder := p.db.Sq.Builder.Select("COUNT(1)").
		From(tableName).
		Where(squirrel.Eq{"user_id": UserId, "product_id": ProductId})

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return false, err
	}

	var count int

	if err = p.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return false, err
	}
	if count != 0 {
		return true, nil
	}
	return false, nil
}

func (p *productRepo) LikeProduct(ctx context.Context, req *entity.Like) (bool, error) {
	data := map[string]any{
		"id":         req.Id,
		"user_id":    req.UserID,
		"product_id": req.ProductID,
	}
	query, args, err := p.db.Sq.Builder.Insert("wishlist").SetMap(data).ToSql()

	if err != nil {
		return false, err
	}

	_, err = p.db.Exec(ctx, query, args...)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *productRepo) DeleteLikeProduct(ctx context.Context, userId, productId string) error {

	sqlStr, args, err := p.db.Sq.Builder.
		Delete("wishlist").
		Where(p.db.Sq.Equal("user_id", userId)).
		Where(p.db.Sq.Equal("product_id", productId)).
		ToSql()

	if err != nil {
		return err
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no sql rows")
	}

	return nil
}

func (u *productRepo) UserWishlist(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error) {
	queryBuilder := u.likesSelectQueryPrefix()
	queryBuilder = queryBuilder.
		From(u.likesTable).
		Where(squirrel.Eq{"user_id": req.Params["user_id"]}).
		OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var like entity.Like
	err = u.db.QueryRow(ctx, query, args[0]).Scan(&like.Id, &like.ProductID, &like.UserID)
	if err != nil {
		return nil, err
	}

	queryBuilder = u.productsSelectQueryPrefix()
	queryBuilder = queryBuilder.From(u.productTable).Where(squirrel.Eq{"id": like.ProductID})

	query, args, err = queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(ctx, query, args[0])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := entity.ListProduct{}
	for rows.Next() {
		var product entity.Product
		if err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Category,
			&product.MadeIn,
			&product.Color,
			&product.Count,
			&product.Cost,
			&product.Discount,
			&product.AgeMin,
			&product.AgeMax,
			&product.ForGender,
			&product.Size,
		); err != nil {
			return nil, err
		}

		products.Products = append(products.Products, &product)
	}

	var count uint64
	total := `SELECT COUNT(user_id) FROM wishlist WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		return nil, err
	}
	products.TotalCount = count

	return &products, nil
}
