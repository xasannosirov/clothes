package v1

import (
	_ "api-gateway/api/docs"
	"api-gateway/api/models"
	pbm "api-gateway/genproto/media_service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/protobuf/encoding/protojson"
)

// Upload photo
// @Summary 		Upload media
// @Description 	Through this api frontent can upload photo and get the link to the media.
// @Tags 			Media
// @Security        BearerAuth
// @Accept 			multipart/form-data
// @Produce         json
// @Param 			productId query string true "productId"
// @Param 			file formData file true "File"
// @Success 		200 {object} models.Response
// @Failure 		500 {object} models.Error
// @Router  		/v1/media/photo [post]
func (h *HandlerV1) UploadMedia(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(7))
	defer cancel()

	endpoint := "13.201.56.179:9000"
	accessKeyID := "abdulaziz"
	secretAccessKey := "abdulaziz"
	bucketName := "clothesstore"
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "BucketAlreadyOwnedByYou" {
		} else {
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

	err = minioClient.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	productId := c.Query("productId")

	file := &models.File{}
	err = c.ShouldBind(&file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	if file.File.Size > 10<<20 {
		c.JSON(http.StatusRequestEntityTooLarge, models.Error{
			Message: "File size cannot be larger than 10 MB",
		})
		return
	}

	ext := filepath.Ext(file.File.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".svg" && ext != ".jpeg"{
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Only .jpg and .png format images are accepted",
		})
		return
	}

	uploadDir := "./media"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	id := uuid.New().String()

	newFilename := id + ext
	uploadPath := filepath.Join(uploadDir, newFilename)

	if err := c.SaveUploadedFile(file.File, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	objectName := newFilename
	contentType := "image/jpeg"
	_, err = minioClient.FPutObject(context.Background(), bucketName, objectName, uploadPath, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	minioURL := fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, objectName)
	_, err = h.Service.MediaService().Create(ctx, &pbm.Media{
		Id:        id,
		ProductId: productId,
		ImageUrl:  minioURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Response: minioURL,
	})
}

// Get Media
// @Summary   Get Media
// @Security  ApiKeyAuth
// @Description Api for getting media by id
// @Tags Media
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.ProductImages
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/media/get/{id} [get]
func (h *HandlerV1) GetMedia(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	response, err := h.Service.MediaService().Get(
		ctx, &pbm.MediaWithProductID{
			ProductId: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete Media
// @Summary Delete Media
// @Security ApiKeyAuth
// @Description Api for delete media
// @Tags Media
// @Accept json
// @Produce json
// @Param id path string true "productId"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/media/delete/{id} [delete]
func (h *HandlerV1) DeleteMedia(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true
	productId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	response, err := h.Service.MediaService().Delete(
		ctx, &pbm.MediaWithProductID{
			ProductId: productId,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Response: response.String(),
	})
}
