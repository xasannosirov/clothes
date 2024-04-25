package services

import (
	"context"
	"time"
	userproto "user-service/genproto/user_service"
	"user-service/internal/entity"
	"user-service/internal/infrastructure/grpc_service_clients"
	"user-service/internal/usecase"

	"go.uber.org/zap"
)

type userRPC struct {
	logger      *zap.Logger
	userUsecase usecase.User
	client      grpc_service_clients.ServiceClients
}

func NewRPC(logger *zap.Logger, userUsecase usecase.User, client *grpc_service_clients.ServiceClients) userproto.UserServiceServer {
	return &userRPC{
		logger:      logger,
		userUsecase: userUsecase,
		client:      *client,
	}
}

func (s userRPC) Create(ctx context.Context, in *userproto.User) (*userproto.UserWithID, error) {
	guid, err := s.userUsecase.Create(ctx, &entity.User{
		GUID:        in.Id,
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Password:    in.Password,
		Gender:      in.Gender,
		Age:         uint8(in.Age),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &userproto.UserWithID{
		Id: guid,
	}, nil
}

func (s userRPC) Update(ctx context.Context, in *userproto.User) (*userproto.User, error) {
	err := s.userUsecase.Update(ctx, &entity.User{
		GUID:        in.Id,
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Password:    in.Password,
		Age:         uint8(in.Age),
		Gender:      in.Gender,
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &userproto.User{
		Id:          in.Id,
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Password:    in.Password,
		Gender:      in.Gender,
		Age:         in.Age,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}, nil
}

func (s userRPC) Delete(ctx context.Context, in *userproto.UserWithID) (*userproto.DeleteUserResponse, error) {
	if err := s.userUsecase.Delete(ctx, in.Id); err != nil {
		s.logger.Error(err.Error())
		return &userproto.DeleteUserResponse{Status: false}, err
	}

	return &userproto.DeleteUserResponse{Status: true}, nil
}

func (s userRPC) Get(ctx context.Context, in *userproto.UserWithID) (*userproto.User, error) {
	user, err := s.userUsecase.Get(ctx, map[string]string{
		"id": in.Id,
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &userproto.User{
		Id:          user.GUID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
		Gender:      user.Gender,
		Age:         int64(user.Age),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s userRPC) GetAll(ctx context.Context, in *userproto.ListUserRequest) (*userproto.ListUserResponse, error) {
	offset := in.Limit * (in.Page - 1)
	users, err := s.userUsecase.List(ctx, uint64(in.Limit), uint64(offset), map[string]string{})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	var response userproto.ListUserResponse
	for _, u := range users {

		temp := &userproto.User{
			Id:          u.GUID,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			Password:    u.Password,
			Gender:      u.Gender,
			Age:         int64(u.Age),
			CreatedAt:   u.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   u.UpdatedAt.Format(time.RFC3339),
		}

		response.Users = append(response.Users, temp)
	}

	return &response, nil
}
