package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"product-service/internal/entity"
	"product-service/internal/pkg/otlp"
	"time"

	"github.com/Masterminds/squirrel"
)

func (u *productRepo) CreateProduct(ctx context.Context, req *entity.Product) (*entity.Product, error) {
	ctx, span := otlp.Start(ctx, "product_grpc-reposiroty", "CreateProduct")
	defer span.End()

	data := map[string]any{
		"id":          req.Id,
		"name":        req.Name,
		"category_id": req.Category,
		"description": req.Description,
		"made_in":     req.MadeIn,
		"color":       req.Color,
		"cost":        req.Cost,
		"count":       req.Count,
		"discount":    req.Discount,
		"age_min":     req.AgeMin,
		"age_max":     req.AgeMax,
		"for_gender":  req.ForGender,
		"size":        req.Size,
	}

	query, args, err := u.db.Sq.Builder.Insert(u.productTable).SetMap(data).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *productRepo) UpdateProduct(ctx context.Context, req *entity.Product) error {
	data := map[string]any{
		"name":        req.Name,
		"description": req.Description,
		"category_id": req.Category,
		"made_in":     req.MadeIn,
		"color":       req.Color,
		"count":       req.Count,
		"cost":        req.Cost,
		"discount":    req.Discount,
		"age_min":     req.AgeMin,
		"age_max":     req.AgeMax,
		"for_gender":  req.ForGender,
		"size":        req.Size,
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.productTable).
		SetMap(data).
		Where(squirrel.Eq{"id": req.Id}).
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

func (u *productRepo) DeleteProduct(ctx context.Context, id string) error {
	caluses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.productTable).
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
		return errors.New("no sql rows")
	}

	return nil
}

func (u *productRepo) GetProduct(ctx context.Context, params map[string]string) (*entity.Product, error) {
	var (
		product entity.Product
	)

	queryBuilder := u.productsSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(squirrel.Eq{key: value})
		}
	}

	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		nullDesc   sql.NullString
		nullColor  sql.NullString
		nullAgeMax sql.NullInt64
	)
	if err = u.db.QueryRow(ctx, query, args...).Scan(
		&product.Id,
		&product.Name,
		&nullDesc,
		&product.Category,
		&product.MadeIn,
		&nullColor,
		&product.Count,
		&product.Cost,
		&product.Discount,
		&product.AgeMin,
		&nullAgeMax,
		&product.ForGender,
		&product.Size,
	); err != nil {
		return nil, err
	}

	if nullDesc.Valid {
		product.Description = nullDesc.String
	}
	if nullColor.Valid {
		product.Color = nullColor.String
	}
	if nullAgeMax.Valid {
		product.AgeMax = nullAgeMax.Int64
	}

	return &product, nil
}

func (u *productRepo) ListProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}
	offset := (req.Page - 1) * req.Limit

	queryBuilder := u.productsSelectQueryPrefix()
	queryBuilder = queryBuilder.Limit(uint64(req.Limit))
	queryBuilder = queryBuilder.Offset(uint64(offset))
	queryBuilder = queryBuilder.OrderBy("created_at")
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		nullDesc   sql.NullString
		nullColor  sql.NullString
		nullAgeMax sql.NullInt64
	)
	for rows.Next() {
		var product entity.Product
		if err = rows.Scan(
			&product.Id,
			&product.Name,
			&nullDesc,
			&product.Category,
			&product.MadeIn,
			&nullColor,
			&product.Count,
			&product.Cost,
			&product.Discount,
			&product.AgeMin,
			&nullAgeMax,
			&product.ForGender,
			&product.Size,
		); err != nil {
			return nil, err
		}

		if nullDesc.Valid {
			product.Description = nullDesc.String
		}
		if nullColor.Valid {
			product.Color = nullColor.String
		}
		if nullAgeMax.Valid {
			product.AgeMax = nullAgeMax.Int64
		}
		products.Products = append(products.Products, &product)
	}

	var count uint64
	total := "SELECT COUNT(*) FROM products WHERE deleted_at IS NULL"
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		products.TotalCount = 0
	}
	products.TotalCount = count

	return products, nil
}

func (u *productRepo) SearchProduct(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}
	offset := req.Limit * (req.Page - 1)
	var (
		nullDesc      sql.NullString
		nullColor     sql.NullString
		nullAgeMin    sql.NullInt64
		nullAgeMax    sql.NullInt64
		nullForGender sql.NullString
		nullSize      sql.NullInt64
	)

	searchName := "'%" + req.Params["name"] + "%'"
	query := fmt.Sprintf("SELECT id, name, description, category_id, made_in, color, count, cost, discount, age_min, age_max, for_gender, size FROM products WHERE deleted_at IS NULL AND name ILIKE %s ORDER BY created_at LIMIT %d OFFSET %d", searchName, req.Limit, offset)

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&nullDesc,
			&product.Category,
			&product.MadeIn,
			&nullColor,
			&product.Count,
			&product.Cost,
			&product.Discount,
			&nullAgeMin,
			&nullAgeMax,
			&nullForGender,
			&nullSize,
		)
		if err != nil {
			return nil, err
		}

		if nullDesc.Valid {
			product.Description = nullDesc.String
		}
		if nullColor.Valid {
			product.Color = nullColor.String
		}
		if nullAgeMin.Valid {
			product.AgeMin = nullAgeMin.Int64
		}
		if nullAgeMax.Valid {
			product.AgeMax = nullAgeMax.Int64
		}
		if nullForGender.Valid {
			product.ForGender = nullForGender.String
		}
		if nullSize.Valid {
			product.Size = nullSize.Int64
		}

		products.Products = append(products.Products, &product)
	}

	var count uint64
	// total := "SELECT COUNT(*) FROM products WHERE deleted_at IS NULL AND name ILIKE $1"
	total := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE deleted_at IS NULL AND name ILIKE %s", searchName)
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		return nil, err
	}
	products.TotalCount = count

	return products, nil
}

func (u *productRepo) GetDisableProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListOrders, error) {
	orders := &entity.ListOrders{}

	queryBuilder := u.ordersSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(squirrel.NotEq{"status": "took"})
	queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"created_at": time.Now().AddDate(0, 0, -7)}).OrderBy("created_at")

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
		var order entity.Order
		if err = rows.Scan(
			&order.Id,
			&order.ProductID,
			&order.UserID,
			&order.Status,
		); err != nil {
			return nil, err
		}

		orders.Orders = append(orders.Orders, &order)
	}

	var count uint64
	total := `SELECT COUNT(*) FROM orders WHERE created_at <= NOW() - INTERVAL '7 days' AND deleted_at IS NULL AND status = 'took'`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		orders.TotalCount = 0
	}
	orders.TotalCount = count

	return orders, nil
}

func (p *productRepo) GetDiscountProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}
	offset := (req.Page - 1) * req.Limit

	queryBuilder := p.productsSelectQueryPrefix()
	queryBuilder = queryBuilder.Where(squirrel.NotEq{"discount": 0}).OrderBy("created_at")
	queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	total := `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL AND discount <> 0`
	if err := p.db.QueryRow(ctx, total).Scan(&count); err != nil {
		products.TotalCount = 0
	}
	products.TotalCount = count

	return products, nil
}
