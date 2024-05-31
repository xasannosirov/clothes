package v1

import (
	_ "api-gateway/api/docs"
	"api-gateway/api/models"
	pbm "api-gateway/genproto/media_service"
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security        BearerAuth
// @Summary 	    Upload media
// @Description     Through this API, frontend can upload a photo and get the link to the media.
// @Tags 			media
// @Accept 			multipart/form-data
// @Produce         json
// @Param 			id query string true "Product ID"
// @Param 			file formData file true "File"
// @Success 		200 {object} string
// @Failure 		500 {object} models.Error
// @Router  		/v1/media/upload-photo [POST]
func (h *HandlerV1) UploadMedia(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	endpoint := h.Config.Minio.Endpoint
	accessKeyID := h.Config.Minio.AccessKeyID
	secretAccessKey := h.Config.Minio.SecretAcessKey
	bucketName := h.Config.Minio.BucketName
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code != "BucketAlreadyOwnedByYou" {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: err.Error(),
			})
			log.Println(err.Error())
			return
		}
	}

	policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "AWS": ["*"]
                },
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::%s/*"]
            }
        ]
    }`, bucketName)

	err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	productId := c.Query("id")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	if file.Size > 35<<20 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "File size cannot be larger than 35 MB",
		})
		return
	}

	ext := filepath.Ext(file.Filename)

	if ext != ".png" && ext != ".jpg" && ext != ".svg" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Only .jpg and .png format images are accepted",
		})
		return
	}

	id := uuid.New().String()
	newFilename := id + ext

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	defer fileContent.Close()

	_, err = minioClient.PutObject(ctx, bucketName, newFilename, fileContent, file.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	minioURL := fmt.Sprintf("https://media.go-clothes.uz/%s/%s", bucketName, newFilename)

	_, err = h.Service.MediaService().Create(ctx, &pbm.Media{
		Id:        id,
		ProductId: productId,
		ImageUrl:  minioURL,
		FileName:  file.Filename,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, minioURL)
}

// @Security  		BearerAuth
// @Summary   		Get Media
// @Description 	Api for getting media by id
// @Tags 			media
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Product ID"
// @Success 		200 {object} models.ProductImages
// @Failure 		404 {object} models.Error
// @Router 			/v1/media/{id} [GET]
func (h *HandlerV1) GetMedia(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	media, err := h.Service.MediaService().Get(
		ctx, &pbm.MediaWithID{
			Id: id,
		})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	if len(media.Images) == 0 {
		c.JSON(http.StatusOK, nil)
		return
	}

	var response models.ProductImages
	for _, image := range media.Images {
		response.Images = append(response.Images, &models.Media{
			Id:        image.Id,
			ProductId: image.ProductId,
			ImageUrl:  image.ImageUrl,
			FileName:  image.FileName,
		})
	}

	c.JSON(http.StatusOK, response)
}

// @Security 		BearerAuth
// @Summary 		Delete Media
// @Description 	Api for delete media
// @Tags 			media
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "productId"
// @Success 		200 {object} string
// @Failure 		404 {object} models.Error
// @Router 			/v1/media/{id} [DELETE]
func (h *HandlerV1) DeleteMedia(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true
	productId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	response, err := h.Service.MediaService().Delete(
		ctx, &pbm.MediaWithID{
			Id: productId,
		})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}
	if !response.Status {
		c.JSON(http.StatusNotFound, models.Error{
			Message: models.NotFoundMessage,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}
