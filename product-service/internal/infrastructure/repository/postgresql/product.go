package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"product-service/internal/entity"
	"product-service/internal/infrastructure/repository"
	"product-service/internal/pkg/otlp"
	"product-service/internal/pkg/postgres"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/k0kubun/pp"
)

const (
	productServiceName = "productService"
	productsTableName  = "products"
	ordersTableName    = "orders"
	commentsTableName  = "comments"
	likesTableName     = "wishlist"
	starsTableName     = "stars"
	savesTableName     = "saves"
)

type productRepo struct {
	productTable string
	orderTable   string
	commentTable string
	likesTable   string
	starsTable   string
	savesTable   string
	db           *postgres.PostgresDB
}

func NewProductsRepo(db *postgres.PostgresDB) repository.Product {
	return &productRepo{
		productTable: productsTableName,
		orderTable:   ordersTableName,
		likesTable:   likesTableName,
		commentTable: commentsTableName,
		starsTable:   starsTableName,
		savesTable:   savesTableName,
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
	)
}

func (u *productRepo) starsSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
		"star",
		"created_at",
		"updated_at",
	)
}

func (u *productRepo) savesSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
		"created_at",
		"updated_at",
	)
}

func (u *productRepo) CreateProduct(ctx context.Context, req *entity.Product) (*entity.Product, error) {
	ctx, span := otlp.Start(ctx, "product_grpc-reposiroty", "CreateProduct")
	defer span.End()

	data := map[string]any{
		"id":          req.Id,
		"name":        req.Name,
		"category":    req.Category,
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

func (u *productRepo) GetProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Product, error) {
	var (
		products []*entity.Product
	)
	queryBuilder := u.productsSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

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
		products = append(products, &product)
	}

	return products, nil
}

func (u *productRepo) UpdateProduct(ctx context.Context, req *entity.Product) error {
	data := map[string]any{
		"name":        req.Name,
		"description": req.Description,
		"category":    req.Category,
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

func (u *productRepo) GetAllOrders(ctx context.Context, req *entity.ListRequest) ([]*entity.Order, error) {
	var (
		orders []*entity.Order
	)

	queryBuilder := u.ordersSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

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

		orders = append(orders, &order)
	}

	return orders, nil
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
	var (
		products []*entity.Product
	)

	queryBuilder := p.productsSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.ILike{"name": req.Name})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.productTable, "SearchProduct"))
	}

	fmt.Println(query)
	rows, err := p.db.Query(ctx, query, args[0])
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

		products = append(products, &product)
	}

	return products, nil
}

func (u *productRepo) RecommentProducts(ctx context.Context, req *entity.Recom) ([]*entity.Product, error) {
	var (
		products []*entity.Product
	)

	queryBuilder := u.productsSelectQueryPrefix()

	queryBuilder = queryBuilder.From(u.productTable).Where(squirrel.Eq{"for_gender": req.Gender})
	queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"age_min": req.Age})
	queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"age_max": req.Age})
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
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

		products = append(products, &product)
	}

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

func (u *productRepo) GetProductComments(ctx context.Context, req *entity.GetWithID) ([]*entity.CommentToProduct, error) {
	var (
		comments []*entity.CommentToProduct
	)

	queryBuilder := u.commentsSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.Eq{"product_id": string(req.ID)}).From(commentsTableName)

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

		comments = append(comments, &comment)
	}
	return comments, nil
}

func (u *productRepo) GetProductLikes(ctx context.Context, req *entity.GetWithID) ([]*entity.LikeProduct, error) {
	var (
		likes []*entity.LikeProduct
	)

	queryBuilder := u.likesSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.Eq{"product_id": string(req.ID)}).From(likesTableName)

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

		likes = append(likes, &like)
	}
	return likes, nil
}

func (u *productRepo) GetProductOrders(ctx context.Context, req *entity.GetWithID) ([]*entity.Order, error) {
	var (
		orders []*entity.Order
	)

	queryBuilder := u.ordersSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.Eq{"product_id": string(req.ID)}).From(ordersTableName)

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

		orders = append(orders, &order)
	}
	return orders, nil
}

func (u *productRepo) GetProductStars(ctx context.Context, req *entity.GetWithID) ([]*entity.StarProduct, error) {
	var (
		stars []*entity.StarProduct
	)

	queryBuilder := u.starsSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.Eq{"product_id": string(req.ID)}).From(starsTableName)

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

		stars = append(stars, &star)
	}
	return stars, nil
}

func (u *productRepo) GetSavedProductsByUserID(ctx context.Context, req string) ([]*entity.Product, error) {
	var (
		products []*entity.Product
		saves    []*entity.SaveProduct
	)

	queryBuilder := u.savesSelectQueryPrefix()

	queryBuilder = queryBuilder.From(u.savesTable).Where(squirrel.Eq{"user_id": req})

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
		saves = append(saves, &save)
	}

	for _, save := range saves {
		queryBuilder = u.productsSelectQueryPrefix()
		queryBuilder = queryBuilder.From(u.productTable).Where(squirrel.Eq{"id": save.ProductID})
		query, args, err := queryBuilder.ToSql()
		if err != nil {
			return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.savesTable, "GetSavedProductByUserID"))
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

			products = append(products, &product)
		}
	}
	return products, nil
}

func (u *productRepo) GetWishlistByUserID(ctx context.Context, req string) ([]*entity.Product, error) {
	var (
		products []*entity.Product
		likes    []*entity.LikeProduct
	)

	queryBuilder := u.likesSelectQueryPrefix()

	queryBuilder = queryBuilder.From(u.likesTable).Where(squirrel.Eq{"user_id": req})

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
		likes = append(likes, &like)
	}

	for _, like := range likes {
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

			products = append(products, &product)
		}
	}
	return products, nil
}

func (u *productRepo) GetOrderedProductsByUserID(ctx context.Context, req string) ([]*entity.Product, error) {
	var (
		products []*entity.Product
		orders   []*entity.Order
	)

	queryBuilder := u.ordersSelectQueryPrefix()

	queryBuilder = queryBuilder.From(u.orderTable).Where(squirrel.Eq{"user_id": req})
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

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
		orders = append(orders, &order)
	}

	for _, order := range orders {
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

			products = append(products, &product)
		}
	}
	return products, nil
}

func (u *productRepo) GetAllComments(ctx context.Context, req *entity.ListRequest) ([]*entity.CommentToProduct, error) {
	var (
		comments []*entity.CommentToProduct
	)
	queryBuilder := u.commentsSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

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
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (u *productRepo) GetAllStars(ctx context.Context, req *entity.ListRequest) ([]*entity.StarProduct, error) {
	var (
		stars []*entity.StarProduct
	)
	queryBuilder := u.starsSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset)).From(u.starsTable)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.starsTable, "GetAllStars"))
	}
	fmt.Println(query)

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
		stars = append(stars, &star)
	}

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

func (u *productRepo) GetDisableProducts(ctx context.Context, req *entity.ListRequest) ([]*entity.Order, error) {
	var (
		orders []*entity.Order
	)

	queryBuilder := u.ordersSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	// queryBuilder = queryBuilder.Where(squirrel.Eq{"status": "ordered"})
	queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"created_at": time.Now().AddDate(0, 0, -10)})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.orderTable, "GetDisableProducts"))
	}
	pp.Println(query)

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

		orders = append(orders, &order)
	}

	return orders, nil
}
