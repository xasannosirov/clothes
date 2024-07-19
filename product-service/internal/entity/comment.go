package entity

import "time"

type Comment struct{
	Id string
	OwnerId string
	ProductId string
	Message string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentUpdateRequest struct {
	Id string
	Message string
	UpdatedAt time.Time
}
type CommentListResponse struct{
	Comment []*Comment
	TotalCount int
}
type DeleteRequest struct {
	Id string 
	Deleted_at time.Time
}

type GetRequest struct{
	Filter map[string]string
}