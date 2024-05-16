package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"user-service/internal/entity"
	"user-service/internal/pkg/otlp"
	"user-service/internal/pkg/postgres"

	"github.com/Masterminds/squirrel"
)

const (
	usersTableName      = "users"
	usersSpanRepoPrefix = "usersRepo"
)

type usersRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewUsersRepo(db *postgres.PostgresDB) *usersRepo {
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
			"role",
			"refresh",
			"created_at",
			"updated_at",
		).From(p.tableName)
}

func (p usersRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {

	ctx, span := otlp.Start(ctx, usersSpanRepoPrefix+"_grpc-reposiroty", "CreateUser")
	defer span.End()

	data := map[string]any{
		"id":           user.GUID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"password":     user.Password,
		"gender":       user.Gender,
		"age":          user.Age,
		"role":         user.Role,
		"refresh":      user.Refresh,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return user, nil
}

func (p usersRepo) Update(ctx context.Context, users *entity.User) (*entity.User, error) {

	ctx, span := otlp.Start(ctx, usersSpanRepoPrefix+"_grpc-reposiroty", "UpdateUser")
	defer span.End()

	clauses := map[string]any{
		"first_name":   users.FirstName,
		"last_name":    users.LastName,
		"email":        users.Email,
		"phone_number": users.PhoneNumber,
		"age":          users.Age,
		"gender":       users.Gender,
		"updated_at":   users.UpdatedAt,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", users.GUID)).
		Where("deleted_at is null").
		ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return nil, p.db.Error(fmt.Errorf("no sql rows"))
	}

	return users, nil
}

func (p usersRepo) Delete(ctx context.Context, guid string) error {
	ctx, span := otlp.Start(ctx, usersSpanRepoPrefix+"_grpc-reposiroty", "DeleteUser")
	defer span.End()

	clauses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", guid)).
		Where("deleted_at IS NULL").
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
	ctx, span := otlp.Start(ctx, usersSpanRepoPrefix+"_grpc-reposiroty", "GetUser")
	defer span.End()

	var (
		user entity.User
	)

	queryBuilder := p.usersSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		} else if key == "email" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		} else if key == "refresh" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
	}
	queryBuilder = queryBuilder.Where("deleted_at is null")
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	var (
		nullPhoneNumber sql.NullString
		nullGender      sql.NullString
		nullAge         sql.NullInt32
		nullRefresh     sql.NullString
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
		&user.Role,
		&nullRefresh,
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
	if nullRefresh.Valid {
		user.Refresh = nullRefresh.String
	}

	return &user, nil
}
func (p usersRepo) GetDelete(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, span := otlp.Start(ctx, usersSpanRepoPrefix+"_grpc-reposiroty", "GetUser")
	defer span.End()

	var (
		user entity.User
	)

	queryBuilder := p.usersSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		} else if key == "email" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		} else if key == "refresh" {
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
		nullRefresh     sql.NullString
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
		&user.Role,
		&nullRefresh,
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
	if nullRefresh.Valid {
		user.Refresh = nullRefresh.String
	}

	return &user, nil
}

func (p usersRepo) List(ctx context.Context, limit uint64, offset uint64, filter map[string]string) ([]*entity.User, error) {

	ctx, span := otlp.Start(ctx, usersSpanRepoPrefix+"_grpc-reposiroty", "ListUsers")
	defer span.End()

	var (
		users []*entity.User
	)
	queryBuilder := p.usersSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	}

	role := filter["role"]
	queryBuilder = queryBuilder.Where(p.db.Sq.Equal("role", role))
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

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
			nullRefresh     sql.NullString
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
			&user.Role,
			&nullRefresh,
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
		if nullRefresh.Valid {
			user.Refresh = nullRefresh.String
		}

		users = append(users, &user)
	}

	return users, nil
}

func (p usersRepo) UniqueEmail(ctx context.Context, request *entity.IsUnique) (*entity.Response, error) {
	ctx, span := otlp.Start(ctx, usersSpanRepoPrefix+"_grpc-reposiroty", "UniqueEmail")
	defer span.End()

	query := `SELECT COUNT(*) FROM users WHERE email = $1 and deleted_at is null`

	var count int
	err := p.db.QueryRow(ctx, query, request.Email).Scan(&count)
	if err != nil {
		return &entity.Response{Status: true}, p.db.Error(err)
	}
	if count != 0 {
		return &entity.Response{Status: true}, nil
	}

	return &entity.Response{Status: false}, nil
}

func (p usersRepo) UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error) {
	ctx, span := otlp.Start(ctx, usersSpanRepoPrefix+"_grpc-reposiroty", "UpdateRefresh")
	defer span.End()

	clauses := map[string]any{
		"refresh": request.RefreshToken,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", request.UserID)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return &entity.Response{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return &entity.Response{Status: false}, p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return &entity.Response{Status: false}, p.db.Error(fmt.Errorf("no sql rows"))
	}

	return &entity.Response{Status: true}, nil
}

func (p usersRepo) UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error) {
	ctx, span := otlp.Start(ctx, usersSpanRepoPrefix+"_grpc-reposiroty", "UpdatePassword")
	defer span.End()

	clauses := map[string]any{
		"password": request.NewPassword,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", request.UserID)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return &entity.Response{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return &entity.Response{Status: false}, p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return &entity.Response{Status: false}, p.db.Error(fmt.Errorf("no sql rows"))
	}

	return &entity.Response{Status: true}, nil
}

func (p usersRepo) Total(ctx context.Context, role string) uint64 {
	query := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND role = $1`

	var count uint64
	if err := p.db.QueryRow(ctx, query, role).Scan(&count); err != nil {
		return 0
	}
	return count
}
