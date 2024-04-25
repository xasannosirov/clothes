package postgres

import (
	"clothes-store/media-service/internal/entity"
	r "clothes-store/media-service/internal/infrastructure/repository"
	configpkg "clothes-store/media-service/internal/pkg/config"
	"clothes-store/media-service/internal/pkg/postgres"
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MediaRepositoryTestSuite struct {
	suite.Suite
	Repository r.MediaStorageI
}

func (s *MediaRepositoryTestSuite) SetupSuite() {
	pgPool, err := postgres.New(configpkg.New())
	if err != nil {
		log.Fatal("Error while connecting database with suite test")
		return
	}

	s.Repository = NewMediaRepo(pgPool)
}

func (s *MediaRepositoryTestSuite) TestSuite() {
	media := &entity.Media{
		Id:         uuid.NewString(),
		Product_Id: "69fdcba8-8d71-47f6-a3a0-977b3469221f",
		Image_Url:  "https://example.com/image.jpg",
	}

	createdMedia, err := s.Repository.CreateMedia(context.Background(), media)

	s.Require().NoError(err)
	s.Require().NotNil(createdMedia)
	s.Equal(createdMedia.Id, media.Id)
	s.Equal(createdMedia.Product_Id, media.Product_Id)
	s.Equal(createdMedia.Image_Url, media.Image_Url)

	//suite test for get all media by product_id
	filter := make(map[string]string)
	filter["product_id"] = createdMedia.Product_Id
	getCreatedMedia, err := s.Repository.GetMediaWithProductId(context.Background(), filter)
	s.Require().NoError(err)
	s.Require().NotNil(getCreatedMedia)

	//suite test for delete media
	id := make(map[string]any)
	id["product_id"] = createdMedia.Product_Id
	err = s.Repository.DeleteMedia(context.Background(), id)
	s.Require().NoError(err)
}

func TestMediaRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MediaRepositoryTestSuite))
}
