package postgresql

import (
	"product-service/internal/infrastructure/repository"
	"product-service/internal/pkg/postgres"

	"github.com/Masterminds/squirrel"
)

const (
	productServiceName = "productService"
	categoryTableName  = "category"
	productsTableName  = "products"
	basketTableName    = "baskets"
	likesTableName     = "wishlist"
	ordersTableName    = "orders"
	commentTableName   = "comments"
)

type productRepo struct {
	categoryTable string
	productTable  string
	basketTable   string
	likesTable    string
	orderTable    string
	commentTable string
	db            *postgres.PostgresDB
}

func NewProductsRepo(db *postgres.PostgresDB) repository.Product {
	return &productRepo{
		productTable:  productsTableName,
		orderTable:    ordersTableName,
		likesTable:    likesTableName,
		categoryTable: categoryTableName,
		basketTable:   basketTableName,
		commentTable: commentTableName,
		db:            db,
	}
}

func (u *productRepo) categorySelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"name",
	).From(u.categoryTable)
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
	).From(u.productTable)
}

func (u *productRepo) ordersSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
		"status",
		"count",
	).From(u.orderTable)
}

func (u *productRepo) likesSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
	).From(u.likesTable)
}

func (u *productRepo) basketsSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"product_id",
		"user_id",
	).From(u.basketTable)
}
func (p *productRepo) comentSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.
		Select(
			"id",
			"product_id",
			"owner_id",
			"message",
			"created_at",
			"updated_at",
		).From(p.commentTable)
}