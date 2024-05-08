package models

import "mime/multipart"

type (
	File struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
	}

	UploadPhotoRes struct {
		URL string `json:"photo_url"`
	}

	MediaResponse struct {
		ErrorCode    int            `json:"error_code"`
		ErrorMessage string         `json:"error_message"`
		Body         UploadPhotoRes `json:"body"`
	}

	Media struct {
		Id        string `json:"id,omitempty"`
		ProductId string `json:"product_id,omitempty"`
		ImageUrl  string `json:"image_url,omitempty"`
		FileName  string `json:"file_name"`
	}

	ProductImages struct {
		Images []*Media `json:"images,omitempty"`
	}
)
