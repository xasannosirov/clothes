package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"user-service/internal/entity"
	"user-service/internal/infrastructure/repository"
	"user-service/internal/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/k0kubun/pp"
)

const (
	usersTableName      = "users"
	usersServiceName    = "userService"
	usersSpanRepoPrefix = "usersRepo"
)

type usersRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewUsersRepo(db *postgres.PostgresDB) repository.Users {
	return &usersRepo{
		tableName: usersTableName,
		db:        db,
	}
}

func (p *usersRepo) usersSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.
		Select(
			"id",
			"first_name",
			"last_name",
			"email",
			"phone_number",
			"password",
			"gender",
			"age",
			"created_at",
			"updated_at",
		).From(p.tableName)
}

func (p usersRepo) Create(ctx context.Context, news *entity.User) error {
	pp.Println(int(news.Age))
	data := map[string]any{
		"id":           news.GUID,
		"first_name":   news.FirstName,
		"last_name":    news.LastName,
		"email":        news.Email,
		"phone_number": news.PhoneNumber,
		"password":     news.Password,
		"gender":       news.Gender,
		"age":          news.Age,
		"created_at":   news.CreatedAt,
		"updated_at":   news.UpdatedAt,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	pp.Println(query)

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		pp.Println(err)
		return p.db.Error(err)
	}

	return nil
}

func (p usersRepo) Update(ctx context.Context, users *entity.User) error {
	clauses := map[string]any{
		"first_name":   users.FirstName,
		"last_name":    users.LastName,
		"email":        users.Email,
		"phone_number": users.PhoneNumber,
		"password":     users.Password,
		"age":          users.Age,
		"gender":       users.Gender,
		"updated_at":   users.UpdatedAt,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", users.GUID)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" update")
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

func (p usersRepo) Delete(ctx context.Context, guid string) error {
	sqlStr, args, err := p.db.Sq.Builder.
		Delete(p.tableName).
		Where(p.db.Sq.Equal("id", guid)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" delete")
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

func (p usersRepo) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	var (
		user entity.User
	)

	queryBuilder := p.usersSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	var (
		nullPhoneNumber sql.NullString
		nullGender      sql.NullString
		nullAge         sql.NullInt32
	)
	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&user.GUID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&nullPhoneNumber,
		&user.Password,
		&nullGender,
		&nullAge,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, p.db.Error(err)
	}
	if nullPhoneNumber.Valid {
		user.PhoneNumber = nullPhoneNumber.String
	}
	if nullGender.Valid {
		user.Gender = nullGender.String
	}
	if nullAge.Valid {
		user.Age = uint8(nullAge.Int32)
	}

	return &user, nil
}

func (p usersRepo) List(ctx context.Context, limit uint64, offset uint64, filter map[string]string) ([]*entity.User, error) {
	var (
		users []*entity.User
	)
	queryBuilder := p.usersSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	users = make([]*entity.User, 0)
	for rows.Next() {
		var (
			user            entity.User
			nullPhoneNumber sql.NullString
			nullGender      sql.NullString
			nullAge         sql.NullInt32
		)
		if err = rows.Scan(
			&user.GUID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&nullPhoneNumber,
			&user.Password,
			&nullGender,
			&nullAge,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}
		if nullPhoneNumber.Valid {
			user.PhoneNumber = nullPhoneNumber.String
		}
		if nullGender.Valid {
			user.Gender = nullGender.String
		}
		if nullAge.Valid {
			user.Age = uint8(nullAge.Int32)
		}

		users = append(users, &user)
	}

	return users, nil
}
