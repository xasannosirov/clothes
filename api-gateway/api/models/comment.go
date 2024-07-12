package models


type Comment struct{
	ID string
	OwnerID string
	ProductID string
	Message string
}

type CommentCreate struct{
	ProductID string
	Message string
}

type CommentUpdate struct{
	ID string
	Message string
}

type  ListComment struct{
	Comment []*Comment
	TotalCount int
}

type Like struct{
	CommentId string
	ProductID string
	OwnerId string
	Status bool
}

type CreateLike struct{
	CommentId string
	ProductID string
	Status bool
}

type Metadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
