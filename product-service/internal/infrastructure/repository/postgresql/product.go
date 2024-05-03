package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"product-service/internal/entity"
	"product-service/internal/infrastructure/repository"
	"product-service/internal/pkg/otlp"
	"product-service/internal/pkg/postgres"

	"github.com/Masterminds/squirrel"
)

const (
	productsTableName    = "products"
	ordersTableName      = "orders"
	productServiceName   = "productService"
	productServicePrefix = "productServiceRepo"
)

type productRepo struct {
	productTable string
	orderTable   string
	db           *postgres.PostgresDB
}

func NewProductsRepo(db *postgres.PostgresDB) repository.Product {
	return &productRepo{
		productTable: productsTableName,
		orderTable:   ordersTableName,
		db:           db,
	}
}

func (u *productRepo) productsSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"name",
		"description",
		"category",
		"made_in",
		"color",
		"count",
		"cost",
		"discount",
		"age_min",
		"age_max",
		"temperature_min",
		"temperature_max",
		"for_gender",
		"size",
		"created_at",
		"updated_at",
	).From(u.productTable)
}

func (u *productRepo) ordersSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
		"status",
		"created_at",
		"updated_at").From(u.orderTable)
}

func (u *productRepo) CreateProduct(ctx context.Context, req *entity.Product) (*entity.Product, error) {
	ctx, span := otlp.Start(ctx, "product_grpc-reposiroty", "CreateProduct")
	defer span.End()

	data := map[string]any{
		"id":              req.Id,
		"name":            req.Name,
		"category":        req.Category,
		"description":     req.Description,
		"made_in":         req.MadeIn,
		"color":           req.Color,
		"cost":            req.Cost,
		"count":           req.Count,
		"discount":        req.Discount,
		"age_min":         req.AgeMin,
		"age_max":         req.AgeMax,
		"temperature_min": req.TemperatureMin,
		"temperature_max": req.TemperatureMax,
		"for_gender":      req.ForGender,
		"size":            req.Size,
		"created_at":      req.CreatedAt,
		"updated_at":      req.UpdatedAt,
	}

	query, args, err := u.db.Sq.Builder.Insert(u.productTable).SetMap(data).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "createProduct"))
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}

	return req, nil
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
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "getProduct"))
	}

	var (
		updatedAt sql.NullTime
		category  sql.NullString
	)
	if err = u.db.QueryRow(ctx, query, args...).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&category,
		&product.MadeIn,
		&product.Color,
		&product.Count,
		&product.Cost,
		&product.Discount,
		&product.AgeMin,
		&product.AgeMax,
		&product.TemperatureMin,
		&product.TemperatureMax,
		&product.ForGender,
		&product.Size,
		&product.CreatedAt,
		&updatedAt,
	); err != nil {
		return nil, u.db.Error(err)
	}

	if updatedAt.Valid {
		product.UpdatedAt = updatedAt.Time
	}
	if category.Valid {
		product.Category = category.String
	}

	return &product, nil
}

func (u *productRepo) GetProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Product, error) {
	var (
		products []*entity.Product
	)
	queryBuilder := u.productsSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "getProducts"))
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()
	var (
		updatedAt sql.NullTime
	)
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
			&product.TemperatureMin,
			&product.TemperatureMax,
			&product.ForGender,
			&product.Size,
			&product.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			product.UpdatedAt = updatedAt.Time
		}
		products = append(products, &product)
	}

	return products, nil
}

func (u *productRepo) UpdateProduct(ctx context.Context, req *entity.Product) error {
	data := map[string]any{
		"name":            req.Name,
		"description":     req.Description,
		"category":        req.Category,
		"made_in":         req.MadeIn,
		"color":           req.Color,
		"count":           req.Count,
		"cost":            req.Cost,
		"discount":        req.Discount,
		"age_min":         req.AgeMin,
		"age_max":         req.AgeMax,
		"temperature_min": req.TemperatureMin,
		"temperature_max": req.TemperatureMax,
		"for_gender":      req.ForGender,
		"size":            req.Size,
		"updated_at":      req.UpdatedAt,
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.productTable).
		SetMap(data).
		Where(squirrel.Eq{"id": req.Id}).
		ToSql()
	fmt.Printf("%v\n\n", sqlStr)
	if err != nil {
		return u.db.ErrSQLBuild(err, u.productTable+" updateProduct")
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return u.db.Error(err)
	}
	if commandTag.RowsAffected() == 0 {
		return u.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (u *productRepo) DeleteProduct(ctx context.Context, id string) error {
	sqlStr, args, err := u.db.Sq.Builder.
		Delete(u.productTable).
		Where(u.db.Sq.Equal("id", id)).
		ToSql()
	if err != nil {
		return u.db.ErrSQLBuild(err, u.productTable+" deleteProduct")
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return u.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return u.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (u *productRepo) CreateOrder(ctx context.Context, req *entity.Order) (*entity.Order, error) {
	data := map[string]any{
		"id":         req.Id,
		"product_id": req.ProductID,
		"user_id":    req.UserID,
		"status":     req.Status,
		"created_at": req.CreatedAt,
		"updated_at": req.UpdatedAt,
	}

	query, args, err := u.db.Sq.Builder.Insert(u.orderTable).SetMap(data).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.orderTable, "CreateOrder"))
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}

	return req, nil
}

func (u *productRepo) GetOrderByID(ctx context.Context, params map[string]string) (*entity.Order, error) {
	var (
		order entity.Order
	)

	queryBuilder := u.ordersSelectQueryPrefix()
	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(squirrel.Eq{key: value})
		}
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.orderTable, "GetOrder"))
	}
	fmt.Printf("%v\n", query)
	var (
		updatedAt sql.NullTime
	)

	if err = u.db.QueryRow(ctx, query, args...).Scan(
		&order.Id,
		&order.ProductID,
		&order.UserID,
		&order.Status,
		&order.CreatedAt,
		&updatedAt,
	); err != nil {
		return nil, u.db.Error(err)
	}

	if updatedAt.Valid {
		order.UpdatedAt = updatedAt.Time
	}

	return &order, nil
}

func (u *productRepo) GetAllOrders(ctx context.Context, req *entity.ListRequest) ([]*entity.Order, error) {
	var (
		orders []*entity.Order
	)

	queryBuilder := u.ordersSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.orderTable, "getAllOrders"))
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	var (
		updatedAt sql.NullTime
	)

	for rows.Next() {
		var order entity.Order
		if err = rows.Scan(
			&order.Id,
			&order.ProductID,
			&order.UserID,
			&order.Status,
			&order.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			order.UpdatedAt = updatedAt.Time
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func (u *productRepo) CancelOrder(ctx context.Context, id string) error {
	sqlStr, args, err := u.db.Sq.Builder.Delete(u.orderTable).Where(u.db.Sq.Equal("id", id)).ToSql()
	if err != nil {
		return u.db.ErrSQLBuild(err, u.orderTable+" CancelOrder")
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return u.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return u.db.Error(fmt.Errorf("no sql rows"))
	}
	return nil
}

func (p *productRepo) GetDiscountProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Product, error) {
	var (
		products []*entity.Product
	)

	queryBuilder := p.productsSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.NotEq{"discount": 0})

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.productTable, "getAllProducts"))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
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
			&product.TemperatureMin,
			&product.TemperatureMax,
			&product.ForGender,
			&product.Size,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}

		products = append(products, &product)
	}

	return products, nil

}
func (p *productRepo) SearchProduct(ctx context.Context, req *entity.Filter) ([]*entity.Product, error) {
	return nil, nil
}

func (p *productRepo) RecommentProducts(ctx context.Context, req *entity.Recom) ([]*entity.Product, error) {
	return nil, nil
}

func (p *productRepo) IsUnique(ctx context.Context, tableName, UserId, ProductId string) (bool, error) {

	queryBuilder := p.db.Sq.Builder.Select("COUNT(1)").
		From(tableName).
		Where(squirrel.Eq{"user_id": UserId, "product_id": ProductId})

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return false, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.productTable, "isUnique"))
	}

	var count int

	if err = p.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return false, p.db.Error(err)

	}
	if count != 0 {
		return true, nil
	}
	return false, nil
}

func (p *productRepo) LikeProduct(ctx context.Context, req *entity.LikeProduct) (bool, error) {
	data := map[string]any{
		"id":         req.Id,
		"user_id":    req.User_id,
		"product_id": req.Product_id,
		"created_at": req.Created_at,
		"updated_at": req.Updated_at,
	}
	query, args, err := p.db.Sq.Builder.Insert("wishlist").SetMap(data).ToSql()

	if err != nil {
		return false, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", "wishlist", "LikeProduct"))
	}

	_, err = p.db.Exec(ctx, query, args...)

	if err != nil {
		return false, p.db.Error(err)
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
		return p.db.ErrSQLBuild(err, "wishlist"+" deleteLikeProduct")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p *productRepo) SaveProduct(ctx context.Context, req *entity.SaveProduct) (bool, error) {
	data := map[string]any{
		"id":         req.Id,
		"user_id":    req.User_id,
		"product_id": req.Product_id,
		"created_at": req.Created_at,
		"updated_at": req.Updated_at,
	}
	query, args, err := p.db.Sq.Builder.Insert("saved").SetMap(data).ToSql()

	if err != nil {
		return false, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", "saved", "LikeProduct"))
	}

	_, err = p.db.Exec(ctx, query, args...)

	if err != nil {
		return false, p.db.Error(err)
	}
	return true, nil
}

func (p *productRepo) DeleteSaveProduct(ctx context.Context, userId, productId string) error {
	sqlStr, args, err := p.db.Sq.Builder.
		Delete("saved").
		Where(p.db.Sq.Equal("user_id", userId)).
		Where(p.db.Sq.Equal("product_id", productId)).
		ToSql()

	if err != nil {
		return p.db.ErrSQLBuild(err, "saved"+" deleteLikeProduct")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p *productRepo) CommentToProduct(ctx context.Context, req *entity.CommentToProduct) (bool, error) {
	data := map[string]any{
		"id":         req.Id,
		"user_id":    req.UserId,
		"product_id": req.Product_Id,
		"comment":    req.Comment,
		"created_at": req.Created_at,
		"updated_at": req.Updated_at,
	}
	query, args, err := p.db.Sq.Builder.Insert("comments").
		SetMap(data).
		ToSql()

	if err != nil {
		return false, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.productTable, "createProduct"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return false, p.db.Error(err)
	}
	return true, nil
}
