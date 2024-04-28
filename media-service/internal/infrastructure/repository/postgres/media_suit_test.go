package postgres

import (
	"context"
	"log"
	"media-service/internal/entity"
	r "media-service/internal/infrastructure/repository"
	configpkg "media-service/internal/pkg/config"
	"media-service/internal/pkg/postgres"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MediaRepositoryTestSuite struct {
	suite.Suite
	Repository r.MediaStorageI
}

// SetupSuite set up envirenment for suite test
func (s *MediaRepositoryTestSuite) SetupSuite() {
	pgPool, err := postgres.New(configpkg.New())
	if err != nil {
		log.Fatal("Error while connecting database with suite test")
		return
	}

	s.Repository = NewMediaRepo(pgPool)
}

// TestSuite test media-service storage methods
func (s *MediaRepositoryTestSuite) TestSuite() {

	// mock info for check media-service storage methods
	defaultProductID := uuid.NewString()
	listMedia := []*entity.Media{
		&entity.Media{
			Id:        uuid.NewString(),
			ProductID: defaultProductID,
			ImageUrl:  "https://clothes-store-management/images/products/1",
		},
		&entity.Media{
			Id:        uuid.NewString(),
			ProductID: defaultProductID,
			ImageUrl:  "https://clothes-store-management/images/products/2",
		},
	}

	// suite test for create media method
	for _, media := range listMedia {
		createdMedia, err := s.Repository.CreateMedia(context.TODO(), media)
		s.Require().NoError(err)
		s.Require().NotNil(createdMedia)
		s.Require().NotEmpty(createdMedia)
	}

	// suite test for get all media by product_id
	filter := make(map[string]string)
	filter["product_id"] = defaultProductID
	listCreatedMedia, err := s.Repository.GetMediaWithProductId(context.Background(), filter)
	s.Require().NoError(err)
	s.Equal(listCreatedMedia[0].ProductID, defaultProductID)
	s.Equal(listCreatedMedia[0].ImageUrl, "https://clothes-store-management/images/products/1")
	s.Require().NotNil(listCreatedMedia[0].Id)
	s.Equal(listCreatedMedia[1].ProductID, defaultProductID)
	s.Equal(listCreatedMedia[1].ImageUrl, "https://clothes-store-management/images/products/2")
	s.Require().NotNil(listCreatedMedia[1].Id)

	// suite test for delete media by product id
	params := make(map[string]any)
	params["product_id"] = defaultProductID
	err = s.Repository.DeleteMedia(context.TODO(), params)
	s.Require().NoError(err)
}

// running media-service suite test
func TestMediaRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MediaRepositoryTestSuite))
}
