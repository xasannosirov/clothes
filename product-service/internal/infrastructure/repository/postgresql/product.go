package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"product-service/internal/entity"
	"product-service/internal/infrastructure/repository"
	"product-service/internal/pkg/otlp"
	"product-service/internal/pkg/postgres"
	"time"

	"github.com/Masterminds/squirrel"
)

const (
	productServiceName = "productService"
	productsTableName  = "products"
	ordersTableName    = "orders"
	commentsTableName  = "comments"
	likesTableName     = "wishlist"
	starsTableName     = "stars"
	savesTableName     = "saves"
	categoryTableName  = "category"
	basketTableName    = "baskets"
)

type productRepo struct {
	productTable  string
	orderTable    string
	commentTable  string
	likesTable    string
	starsTable    string
	savesTable    string
	categoryTable string
	basketTable   string
	db            *postgres.PostgresDB
}

func NewProductsRepo(db *postgres.PostgresDB) repository.Product {
	return &productRepo{
		productTable:  productsTableName,
		orderTable:    ordersTableName,
		likesTable:    likesTableName,
		commentTable:  commentsTableName,
		starsTable:    starsTableName,
		savesTable:    savesTableName,
		categoryTable: categoryTableName,
		basketTable:   basketTableName,
		db:            db,
	}
}

func (u *productRepo) productsSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"name",
		"description",
		"category_id",
		"made_in",
		"color",
		"count",
		"cost",
		"discount",
		"age_min",
		"age_max",
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
		"updated_at",
	).From(u.orderTable)
}

func (u *productRepo) commentsSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
		"comment",
		"created_at",
		"updated_at").From(u.commentTable)
}

func (u *productRepo) likesSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
		"created_at",
		"updated_at",
	).From(u.likesTable)
}

func (u *productRepo) starsSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
		"star",
		"created_at",
		"updated_at",
	).From(u.starsTable)
}

func (u *productRepo) savesSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
		"created_at",
		"updated_at",
	).From(u.savesTable)
}

func (u *productRepo) categorySelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"name",
		"created_at",
		"updated_at",
	).From(u.categoryTable)
}

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
		"created_at":  req.CreatedAt,
		"updated_at":  req.UpdatedAt,
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

	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "getProduct"))
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
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, u.db.Error(err)
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
func (u *productRepo) GetProductDelete(ctx context.Context, params map[string]string) (*entity.Product, error) {
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
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, u.db.Error(err)
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

func (u *productRepo) GetProducts(ctx context.Context, req *entity.ListProductRequest) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}
	queryBuilder := u.productsSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Name != "SkottAdkins" {
		queryBuilder = queryBuilder.Where("name ILIKE " + "'%" + req.Name + "%'")
	}
	queryBuilder = queryBuilder.Where(" deleted_at IS NULL LIMIT $1 OFFSET $2").OrderBy("created_at")

	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "getProducts"))
	}

	rows, err := u.db.Query(ctx, query, req.Limit, offset)
	if err != nil {
		return nil, u.db.Error(err)
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
			&product.CreatedAt,
			&product.UpdatedAt,
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
	if req.Name != "SkottAdkins" {
		total += " AND name ILIKE " + "'%" + req.Name + "%'"
	}
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		products.TotalCount = 0
	}
	products.TotalCount = count

	return products, nil
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
		"updated_at":  req.UpdatedAt,
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.productTable).
		SetMap(data).
		Where(squirrel.Eq{"id": req.Id}).
		ToSql()
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
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.orderTable, "GetOrder"))
	}

	if err = u.db.QueryRow(ctx, query, args...).Scan(
		&order.Id,
		&order.ProductID,
		&order.UserID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return nil, u.db.Error(err)
	}

	return &order, nil
}

func (u *productRepo) CancelOrder(ctx context.Context, id string) error {
	clauses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}
	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.orderTable).
		SetMap(clauses).
		Where(u.db.Sq.Equal("id", id)).
		Where("deleted_at IS NULL").
		ToSql()
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

func (u *productRepo) GetAllOrders(ctx context.Context, req *entity.ListRequest) (*entity.ListOrders, error) {
	orders := &entity.ListOrders{}

	queryBuilder := u.ordersSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	queryBuilder = queryBuilder.Where("deleted_at IS NULL").OrderBy("created_at").OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.orderTable, "getAllOrders"))
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
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}

		orders.Orders = append(orders.Orders, &order)
	}

	var count uint64
	total := `SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		orders.TotalCount = 0
	}
	orders.TotalCount = count

	return orders, nil
}

func (p *productRepo) GetDiscountProducts(ctx context.Context, req *entity.ListRequest) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}

	queryBuilder := p.productsSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.NotEq{"discount": 0}).OrderBy("created_at")

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
			&product.ForGender,
			&product.Size,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
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

func (p *productRepo) SearchProduct(ctx context.Context, req *entity.Filter) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}

	queryBuilder := p.productsSelectQueryPrefix()

	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.productTable, "SearchProduct"))
	}

	name := "%" + req.Name + "%"
	query += " WHERE name ILIKE $1"

	rows, err := p.db.Query(ctx, query, name)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	var (
		nullDesc      sql.NullString
		nullColor     sql.NullString
		nullAgeMin    sql.NullInt64
		nullAgeMax    sql.NullInt64
		nullForGender sql.NullString
		nullSize      sql.NullInt64
		updatedAt     sql.NullTime
	)

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
			&product.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, p.db.Error(err)
		}

		if updatedAt.Valid {
			product.UpdatedAt = updatedAt.Time
		}

		products.Products = append(products.Products, &product)
	}

	return products, nil
}

func (u *productRepo) RecommentProducts(ctx context.Context, req *entity.Recom) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}

	queryBuilder := u.productsSelectQueryPrefix()

	queryBuilder = queryBuilder.From(u.productTable).Where(squirrel.Eq{"for_gender": req.Gender})
	queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"age_min": req.Age})
	queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"age_max": req.Age})
	queryBuilder = queryBuilder.Where("deleted_at IS NULL").OrderBy("created_at")
	queryBuilder = queryBuilder.Offset(0).Limit(10)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "Recommendations"))
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	var updatedAt sql.NullTime

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
			&product.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			product.UpdatedAt = updatedAt.Time
		}

		products.Products = append(products.Products, &product)
	}

	var count uint64
	total := `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL AND for_gener = $1 AND age_min <= $2 AND age_max >= $3`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		products.TotalCount = 0
	}
	products.TotalCount = count

	return products, nil
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
		"user_id":    req.UserID,
		"product_id": req.ProductID,
		"created_at": req.CreatedAt,
		"updated_at": req.UpdatedAt,
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
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p *productRepo) SaveProduct(ctx context.Context, req *entity.SaveProduct) (bool, error) {
	data := map[string]any{
		"id":         req.Id,
		"user_id":    req.UserID,
		"product_id": req.ProductID,
		"created_at": req.CreatedAt,
		"updated_at": req.UpdatedAt,
	}
	query, args, err := p.db.Sq.Builder.Insert("saves").SetMap(data).ToSql()

	if err != nil {
		return false, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", "saves", "SaveProduct"))
	}

	_, err = p.db.Exec(ctx, query, args...)

	if err != nil {
		return false, p.db.Error(err)
	}
	return true, nil
}

func (p *productRepo) DeleteSaveProduct(ctx context.Context, userId, productId string) error {
	sqlStr, args, err := p.db.Sq.Builder.
		Delete("saves").
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
		"user_id":    req.UserID,
		"product_id": req.ProductID,
		"comment":    req.Comment,
		"created_at": req.CreatedAt,
		"updated_at": req.UpdatedAt,
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

func (u *productRepo) GetProductComments(ctx context.Context, req *entity.GetWithID) (*entity.ListComments, error) {
	comments := &entity.ListComments{}

	queryBuilder := u.commentsSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.Eq{"product_id": string(req.ID)}).From(commentsTableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL").OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.commentTable, "GetProductComments"))
	}

	rows, err := u.db.Query(ctx, query, args[0])
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	var (
		updatedAt sql.NullTime
	)
	for rows.Next() {
		var comment entity.CommentToProduct
		if err = rows.Scan(
			&comment.Id,
			&comment.ProductID,
			&comment.UserID,
			&comment.Comment,
			&comment.CreatedAt,
			&updatedAt); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			comment.UpdatedAt = updatedAt.Time
		}

		comments.Comments = append(comments.Comments, &comment)
	}

	var count uint64
	total := `SELECT COUNT(product_id) FROM comments WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		comments.TotalCount = 0
	}
	comments.TotalCount = count

	return comments, nil
}

func (u *productRepo) GetProductLikes(ctx context.Context, req *entity.GetWithID) (*entity.ListLikes, error) {
	likes := &entity.ListLikes{}

	queryBuilder := u.likesSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.Eq{"product_id": string(req.ID)}).OrderBy("created_at").From(likesTableName)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.likesTable, "GetProductLikes"))
	}

	rows, err := u.db.Query(ctx, query, args[0])
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	var (
		updatedAt sql.NullTime
	)
	for rows.Next() {
		var like entity.LikeProduct
		if err = rows.Scan(
			&like.Id,
			&like.ProductID,
			&like.UserID,
			&like.CreatedAt,
			&updatedAt); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			like.UpdatedAt = updatedAt.Time
		}

		likes.Likes = append(likes.Likes, &like)
	}

	var count uint64
	total := `SELECT COUNT(product_id) FROM wishlist WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		likes.TotalCount = 0
	}
	likes.TotalCount = count

	return likes, nil
}

func (u *productRepo) GetProductOrders(ctx context.Context, req *entity.GetWithID) (*entity.ListOrders, error) {
	orders := &entity.ListOrders{}

	queryBuilder := u.ordersSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.Eq{"product_id": string(req.ID)}).OrderBy("created_at").From(ordersTableName)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.orderTable, "GetProductOrders"))
	}

	rows, err := u.db.Query(ctx, query, args[0])
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
			&updatedAt); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			order.UpdatedAt = updatedAt.Time
		}

		orders.Orders = append(orders.Orders, &order)
	}

	var count uint64
	total := `SELECT COUNT(product_id) FROM orders WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		orders.TotalCount = 0
	}
	orders.TotalCount = count

	return orders, nil
}

func (u *productRepo) GetProductStars(ctx context.Context, req *entity.GetWithID) (*entity.ListStars, error) {
	stars := &entity.ListStars{}

	queryBuilder := u.starsSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.Eq{"product_id": string(req.ID)}).OrderBy("created_at").From(starsTableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.starsTable, "GetProductStars"))
	}

	rows, err := u.db.Query(ctx, query, args[0])
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	var (
		updatedAt sql.NullTime
	)
	for rows.Next() {
		var star entity.StarProduct
		if err = rows.Scan(
			&star.Id,
			&star.ProductID,
			&star.UserID,
			&star.Stars,
			&star.CreatedAt,
			&updatedAt); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			star.UpdatedAt = updatedAt.Time
		}

		stars.Stars = append(stars.Stars, &star)
	}

	var count uint64
	total := `SELECT COUNT(product_id) FROM stars WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		stars.TotalCount = 0
	}
	stars.TotalCount = count

	return stars, nil
}

func (u *productRepo) GetSavedProductsByUserID(ctx context.Context, req string) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}
	saves := &entity.ListSaves{}

	queryBuilder := u.savesSelectQueryPrefix()

	queryBuilder = queryBuilder.From(u.savesTable).Where(squirrel.Eq{"user_id": req})
	queryBuilder = queryBuilder.Where("deleted_at IS NULL").OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.savesTable, "GetSavedProductByUserID"))
	}

	rows, err := u.db.Query(ctx, query, args[0])
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	var updatedAt sql.NullTime

	for rows.Next() {
		var save entity.SaveProduct
		if err = rows.Scan(
			&save.Id,
			&save.ProductID,
			&save.UserID,
			&save.CreatedAt,
			&save.UpdatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}
		saves.Saves = append(saves.Saves, &save)
	}

	for _, save := range saves.Saves {
		queryBuilder = u.productsSelectQueryPrefix()
		queryBuilder = queryBuilder.From(u.productTable).Where(squirrel.Eq{"id": save.ProductID})
		query, args, err := queryBuilder.ToSql()
		if err != nil {
			log.Println("Saved product deleted", save.ProductID)
			continue
		}
		rows, err := u.db.Query(ctx, query, args[0])
		if err != nil {
			return nil, u.db.Error(err)
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
				&product.CreatedAt,
				&updatedAt,
			); err != nil {
				return nil, u.db.Error(err)
			}

			if updatedAt.Valid {
				product.UpdatedAt = updatedAt.Time
			}

			products.Products = append(products.Products, &product)
		}
	}

	var count uint64
	total := `SELECT COUNT(user_id) FROM saves WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		products.TotalCount = 0
	}
	products.TotalCount = count

	return products, nil
}

func (u *productRepo) GetWishlistByUserID(ctx context.Context, req string) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}
	likes := &entity.ListLikes{}

	queryBuilder := u.likesSelectQueryPrefix()

	queryBuilder = queryBuilder.From(u.likesTable).Where(squirrel.Eq{"user_id": req}).OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.savesTable, "GetWishlistByUserID"))
	}

	rows, err := u.db.Query(ctx, query, args[0])
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	var updatedAt sql.NullTime

	for rows.Next() {
		var like entity.LikeProduct
		if err = rows.Scan(
			&like.Id,
			&like.ProductID,
			&like.UserID,
			&like.CreatedAt,
			&like.UpdatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}
		likes.Likes = append(likes.Likes, &like)
	}

	for _, like := range likes.Likes {
		queryBuilder = u.productsSelectQueryPrefix()
		queryBuilder = queryBuilder.From(u.productTable).Where(squirrel.Eq{"id": like.ProductID})
		query, args, err := queryBuilder.ToSql()
		if err != nil {
			return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.savesTable, "GetWishlistByUserID"))
		}
		rows, err := u.db.Query(ctx, query, args[0])
		if err != nil {
			return nil, u.db.Error(err)
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
				&product.CreatedAt,
				&updatedAt,
			); err != nil {
				return nil, u.db.Error(err)
			}

			if updatedAt.Valid {
				product.UpdatedAt = updatedAt.Time
			}

			products.Products = append(products.Products, &product)
		}
	}

	var count uint64
	total := `SELECT COUNT(user_id) FROM wishlist WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		products.TotalCount = 0
	}
	products.TotalCount = count

	return products, nil
}

func (u *productRepo) GetOrderedProductsByUserID(ctx context.Context, req string) (*entity.ListProduct, error) {
	products := &entity.ListProduct{}
	orders := &entity.ListOrders{}

	queryBuilder := u.ordersSelectQueryPrefix()

	queryBuilder = queryBuilder.From(u.orderTable).Where(squirrel.Eq{"user_id": req})
	queryBuilder = queryBuilder.Where("deleted_at IS NULL").OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.savesTable, "GetOrderProductByUserID"))
	}

	rows, err := u.db.Query(ctx, query, args[0])
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	var updatedAt sql.NullTime

	for rows.Next() {
		var order entity.Order
		if err = rows.Scan(
			&order.Id,
			&order.ProductID,
			&order.UserID,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}
		orders.Orders = append(orders.Orders, &order)
	}

	for _, order := range orders.Orders {
		queryBuilder = u.productsSelectQueryPrefix()
		queryBuilder = queryBuilder.From(u.productTable).Where(squirrel.Eq{"id": order.ProductID})
		query, args, err := queryBuilder.ToSql()
		if err != nil {
			return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.orderTable, "GetOrderProductsByUserID"))
		}
		rows, err := u.db.Query(ctx, query, args[0])
		if err != nil {
			return nil, u.db.Error(err)
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
				&product.CreatedAt,
				&updatedAt,
			); err != nil {
				return nil, u.db.Error(err)
			}

			if updatedAt.Valid {
				product.UpdatedAt = updatedAt.Time
			}

			products.Products = append(products.Products, &product)
		}
	}

	var count uint64
	total := `SELECT COUNT(user_id) FROM orders WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		products.TotalCount = 0
	}
	products.TotalCount = count

	return products, nil
}

func (u *productRepo) GetAllComments(ctx context.Context, req *entity.ListRequest) (*entity.ListComments, error) {
	comments := &entity.ListComments{}
	queryBuilder := u.commentsSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit
	queryBuilder = queryBuilder.OrderBy("created_at")
	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.commentTable, "GetAllComments"))
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
		var comment entity.CommentToProduct
		if err = rows.Scan(
			&comment.Id,
			&comment.ProductID,
			&comment.UserID,
			&comment.Comment,
			&comment.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			comment.UpdatedAt = updatedAt.Time
		}
		comments.Comments = append(comments.Comments, &comment)
	}

	var count uint64
	total := `SELECT COUNT(*) FROM comments WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		comments.TotalCount = 0
	}
	comments.TotalCount = count

	return comments, nil
}

func (u *productRepo) GetAllStars(ctx context.Context, req *entity.ListRequest) (*entity.ListStars, error) {
	stars := &entity.ListStars{}
	queryBuilder := u.starsSelectQueryPrefix()
	queryBuilder = queryBuilder.OrderBy("created_at")

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset)).From(u.starsTable)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.starsTable, "GetAllStars"))
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
		var star entity.StarProduct
		if err = rows.Scan(
			&star.Id,
			&star.ProductID,
			&star.UserID,
			&star.Stars,
			&star.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			star.UpdatedAt = updatedAt.Time
		}
		stars.Stars = append(stars.Stars, &star)
	}

	var count uint64
	total := `SELECT COUNT(*) FROM stars WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		stars.TotalCount = 0
	}
	stars.TotalCount = count

	return stars, nil
}

func (u *productRepo) StarProduct(ctx context.Context, req *entity.StarProduct) (*entity.StarProduct, error) {
	data := map[string]any{
		"id":         req.Id,
		"product_id": req.ProductID,
		"user_id":    req.UserID,
		"star":       req.Stars,
		"created_at": req.CreatedAt,
		"updated_at": req.UpdatedAt,
	}

	query, args, err := u.db.Sq.Builder.Insert(u.starsTable).SetMap(data).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.starsTable, "StarProduct"))
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}

	return req, nil
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
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.orderTable, "GetDisableProducts"))
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		if err = rows.Scan(
			&order.Id,
			&order.ProductID,
			&order.UserID,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, u.db.Error(err)
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

func (u *productRepo) CreateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error) {
	ctx, span := otlp.Start(ctx, "product_grpc-reposiroty", "CreateCategory")
	defer span.End()

	data := map[string]any{
		"id":         req.ID,
		"name":       req.Name,
		"created_at": time.Now().Format(time.RFC3339),
		"updated_at": time.Now().Format(time.RFC3339),
	}

	query, args, err := u.db.Sq.Builder.Insert(u.categoryTable).SetMap(data).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.categoryTable, "CreateCategory"))
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}

	return req, nil
}

func (u *productRepo) DeleteCategory(ctx context.Context, id string) error {
	caluses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.categoryTable).
		SetMap(caluses).
		Where(u.db.Sq.Equal("id", id)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return u.db.ErrSQLBuild(err, u.categoryTable+" deleteCategory")
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

func (u *productRepo) ListCategory(ctx context.Context, req *entity.ListRequest) (*entity.LiestCategory, error) {
	categories := &entity.LiestCategory{}
	queryBuilder := u.categorySelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	queryBuilder = queryBuilder.Where("deleted_at IS NULL").OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "getCategories"))
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var category entity.Category
		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}

		categories.Categories = append(categories.Categories, &category)
	}

	var count uint64
	totol := `SELECT COUNT(*) FROM category WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, totol).Scan(&count); err != nil {
		categories.TotalCount = 0
	}
	categories.TotalCount = count

	return categories, nil
}

func (u *productRepo) GetCategory(ctx context.Context, id string) (*entity.Category, error) {

	queryBuilder := u.categorySelectQueryPrefix()

	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(u.db.Sq.Equal("id", id))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "getCategories"))
	}

	row := u.db.QueryRow(ctx, query, args...)

	var category entity.Category
	if err = row.Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &category, nil
}

func (u *productRepo) UpdateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error) {
	data := map[string]any{
		"name":       req.Name,
		"updated_at": req.UpdatedAt,
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.categoryTable).
		SetMap(data).
		Where(squirrel.Eq{"id": req.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	if commandTag.RowsAffected() == 0 {
		return nil, err
	}

	return req, nil
}
