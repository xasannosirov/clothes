package postgresql

import (
	"context"
	"fmt"
	"product-service/internal/entity"
)

func (p *productRepo) CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	data := map[string]any{
		"id":         comment.Id,
		"product_id": comment.ProductId,
		"owner_id":   comment.OwnerId,
		"message":    comment.Message,
		"created_at": comment.CreatedAt,
		"updated_at": comment.UpdatedAt,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.commentTable).SetMap(data).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.commentTable, "create"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return comment, nil
}

func (p *productRepo) UpdateComment(ctx context.Context, comment *entity.CommentUpdateRequest) (*entity.Comment, error) {
	clauses := map[string]any{
		"message":    comment.Message,
		"updated_at": comment.UpdatedAt,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.commentTable).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", comment.Id)).
		Where("deleted_at is null").
		ToSql()

	if err != nil {
		return nil, p.db.ErrSQLBuild(err, p.commentTable+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return nil, p.db.Error(fmt.Errorf("no sql rows"))
	}

	filter := map[string]string{
		"id" : comment.Id,
	}
	return p.GetComment(ctx, &entity.GetRequest{
		Filter: filter,
	})
}

func (p *productRepo) DeleteComment(ctx context.Context, req *entity.DeleteRequest) error {
	clauses := map[string]interface{}{
		"deleted_at": req.Deleted_at,
	}

	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.commentTable).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", req.Id)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.commentTable+" delete")
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

func (p *productRepo) GetComment(ctx context.Context, req *entity.GetRequest) (*entity.Comment, error) {
	var (
		comment entity.Comment
		cnt     int
	)

	queryBuilder := p.comentSelectQueryPrefix()

	for key, value := range req.Filter {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		} else if key == "del" {
			cnt++
		}
	}
	if cnt == 0 {
		queryBuilder = queryBuilder.Where("deleted_at is null")
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.commentTable, "get"))
	}

	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&comment.Id,
		&comment.ProductId,
		&comment.OwnerId,
		&comment.Message,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	); err != nil {
		return nil, p.db.Error(err)
	}

	return &comment, nil
}

func (p *productRepo) ListComment(ctx context.Context, req *entity.ListRequest) (*entity.CommentListResponse, error) {
	var (
		comments entity.CommentListResponse
	)
	offset := (req.Page - 1) * req.Limit
	queryBuilder := p.comentSelectQueryPrefix()

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	for key, value := range req.Filter {
		if key == "owner_id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
		if key == "post_id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
	}

	queryBuilder = queryBuilder.Where("deleted_at IS NULL").OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.commentTable, "list"))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment entity.Comment

		if err = rows.Scan(
			&comment.Id,
			&comment.ProductId,
			&comment.OwnerId,
			&comment.Message,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}

		comments.Comment = append(comments.Comment, &comment)
	}

	var count uint64

	queryBuilder = p.db.Sq.Builder.Select("COUNT(*)").
		From(p.commentTable).
		Where("deleted_at is null")

	query, _, err = queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.commentTable, "list"))
	}

	if err := p.db.QueryRow(ctx, query).Scan(&count); err != nil {
		comments.TotalCount = 0
	}
	comments.TotalCount = int(count)

	return &comments, nil
}
