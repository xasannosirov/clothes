package v1

import (
	"api-gateway/api/models"
	pb "api-gateway/genproto/product_service"
	"api-gateway/internal/pkg/regtool"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Security      BearerAuth
// @Summary  	  Create Comment
// @Description   This api for create commment to post
// @Tags   		  comment
// @Accept 	      json
// @Produce 	  json
// @Param 		  comment body models.CommentCreate true "Comment Create Model"
// @Succes        201  {object} models.CreateResponse
// @Failure       401 {object} models.Error
// @Failure       403 {object} models.Error
// @Failure       500 {object} models.Error
// @Router        /v1/comment  [POST]
func (h *HandlerV1) CreateComment(c *gin.Context) {
	var (
		body        models.CommentCreate
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

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

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(500, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	userId, statusCode := regtool.GetIdFromToken(c.Request, &h.Config)
	if statusCode != 0 {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "oops something went wrong",
		})
	}
	Comment, err := h.Service.ProductService().CreateComment(ctx, &pb.Comment{
		OwnerId:   userId,
		ProductId: body.ProductID,
		Message:   body.Message,
	})
	if err != nil {
		c.JSON(401, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusCreated, models.CreateResponse{
		ID: Comment.Id,
	})
}

// @Security  		BearerAuth
// @Summary   		Update Comment
// @Description 	Api for update a Comment
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			comment body models.CommentUpdate true "Update Comment Model"
// @Success 		200 {object} models.Comment
// @Failure 		400 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/comment [PUT]
func (h *HandlerV1) UpdateComment(c *gin.Context) {
	var (
		body        models.CommentUpdate
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

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

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	comment, err := h.Service.ProductService().UpdateComment(ctx, &pb.CommentUpdateRequst{
		Id:      body.ID,
		Message: body.Message,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Comment{
		ID:        comment.Id,
		OwnerID:   comment.OwnerId,
		ProductID: comment.ProductId,
		Message:   comment.Message,
	})
}

// @Security  		BearerAuth
// @Summary   		Delete Comment
// @Description 	Api for delete a comment
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Comment ID"
// @Success 		200 {object} bool
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/comment/{id} [DELETE]
func (h *HandlerV1) DeleteComment(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

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

	userID := c.Param("id")

	_, err = h.Service.ProductService().DeleteComment(ctx, &pb.CommentDeleteRequest{
		Id: userID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, true)
}

// @Security  		BearerAuth
// @Summary   		Get Comment
// @Description 	Api for getting a comment
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Comment ID"
// @Success 		200 {object} models.Comment
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/comment/{id} [GET]
func (h *HandlerV1) GetComment(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

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

	userID := c.Param("id")

	filter := map[string]string{
		"id": userID,
	}
	comment, err := h.Service.ProductService().GetComment(ctx, &pb.CommentGetRequst{
		Filter: filter,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.Comment{
		ID:        userID,
		OwnerID:   comment.OwnerId,
		ProductID: comment.ProductId,
		Message:   comment.Message,
	})
}

// @Security  		BearerAuth
// @Summary   		List Comment
// @Description 	Api for getting list comment
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			page query uint64 true "Page"
// @Param 			limit query uint64 true "Limit"
// @Success 		200 {object} models.ListComment
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/comments [GET]
func (h *HandlerV1) ListComment(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

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

	page := c.Query("page")
	limit := c.Query("limit")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	filter := map[string]string{}
	listComment, err := h.Service.ProductService().ListComment(ctx, &pb.CommentListRequest{
		Page:   int64(pageInt),
		Limit:  int64(limitInt),
		Filter: filter,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var comments []*models.Comment
	for _, comment := range listComment.Comments {
		comments = append(comments, &models.Comment{
			ID:        comment.Id,
			ProductID: comment.ProductId,
			OwnerID:   comment.OwnerId,
			Message:   comment.Message,
		})
	}

	c.JSON(http.StatusOK, models.ListComment{
		Comment:    comments,
		TotalCount: int(listComment.TotalCount),
	})
}

// @Security  		BearerAuth
// @Summary   		List Comment
// @Description 	Api for getting post's comment
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			page query int true "Page"
// @Param 			limit query int true "Limit"
// @Param 			id query string true "User Id"
// @Success 		200 {object} models.ListComment
// @Failure 		404 {object} models.Error
// @Failure 		401 {object} models.Error
// @Failure 		403 {object} models.Error
// @Failure 		500 {object} models.Error
// @Router 			/v1/post/comments [GET]
func (h *HandlerV1) GetAllCommentByPostId(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

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

	UserId := c.Query("id")
	page := c.Query("page")
	limit := c.Query("limit")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	filter := map[string]string{
		"product_id": UserId,
	}
	listComment, err := h.Service.ProductService().ListComment(ctx, &pb.CommentListRequest{
		Page:   int64(pageInt),
		Limit:  int64(limitInt),
		Filter: filter,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Message: err.Error(),
		})
		log.Println(err.Error())
		return
	}

	var comments []*models.Comment
	for _, comment := range listComment.Comments {
		comments = append(comments, &models.Comment{
			ID:        comment.Id,
			OwnerID:   comment.OwnerId,
			ProductID: comment.ProductId,
			Message:   comment.Message,
		})
	}

	c.JSON(http.StatusOK, models.ListComment{
		Comment:    comments,
		TotalCount: int(listComment.TotalCount),
	})
}
