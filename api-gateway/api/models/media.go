package models

import "mime/multipart"

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UploadPhotoRes struct {
	URL string `json:"photo_url"`
}

type MediaResponse struct {
	ErrorCode    int            `json:"error_code"`
	ErrorMessage string         `json:"error_message"`
	Body         UploadPhotoRes `json:"body"`
}
type Media struct{
	Id                   string   `json:"id,omitempty"`
	ProductId            string   `json:"product_id,omitempty"`
	ImageUrl             string   `json:"image_url,omitempty"`	
}

type ProductImages struct {
	Images               []*Media `json:"images,omitempty"`
}
