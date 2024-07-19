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

func (s userRPC) CreateUser(ctx context.Context, in *userproto.User) (*userproto.UserWithGUID, error) {
	user, err := s.userUsecase.Create(ctx, &entity.User{
		GUID:        in.Id,
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Password:    in.Password,
		Gender:      in.Gender,
		Age:         uint8(in.Age),
		Role:        in.Role,
		Refresh:     in.Refresh,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &userproto.UserWithGUID{
		Guid: user.GUID,
	}, nil
}

func (s userRPC) UpdateUser(ctx context.Context, in *userproto.User) (*userproto.User, error) {
	user, err := s.userUsecase.Update(ctx, &entity.User{
		GUID:        in.Id,
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		Age:         uint8(in.Age),
		Gender:      in.Gender,
		UpdatedAt:   time.Now().UTC(),
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

func (s userRPC) DeleteUser(ctx context.Context, in *userproto.UserWithGUID) (*userproto.ResponseStatus, error) {
	if err := s.userUsecase.Delete(ctx, in.Guid); err != nil {
		s.logger.Error(err.Error())
		return &userproto.ResponseStatus{Status: false}, err
	}

	return &userproto.ResponseStatus{Status: true}, nil
}

func (s userRPC) GetUser(ctx context.Context, in *userproto.Filter) (*userproto.User, error) {
	user, err := s.userUsecase.Get(ctx, in.Filter)

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
		Role:        user.Role,
		Refresh:     user.Refresh,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s userRPC) GetUserDelete(ctx context.Context, in *userproto.Filter) (*userproto.User, error) {
	user, err := s.userUsecase.GetDelete(ctx, in.Filter)

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
		Role:        user.Role,
		Refresh:     user.Refresh,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s userRPC) GetAllUsers(ctx context.Context, in *userproto.ListUserRequest) (*userproto.ListUserResponse, error) {
	offset := in.Limit * (in.Page - 1)
	users, err := s.userUsecase.List(ctx, uint64(in.Limit), uint64(offset), map[string]string{
		"role": in.Role,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	response := userproto.ListUserResponse{}
	for _, u := range users {
		temp := &userproto.User{
			Id:          u.GUID,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			Password:    u.Password,
			Gender:      u.Gender,
			Role:        u.Role,
			Refresh:     u.Refresh,
			Age:         int64(u.Age),
			CreatedAt:   u.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   u.UpdatedAt.Format(time.RFC3339),
		}

		response.Users = append(response.Users, temp)
	}
	response.TotalCount = s.userUsecase.Total(ctx, in.Role)

	return &response, nil
}

func (s userRPC) UniqueEmail(ctx context.Context, in *userproto.IsUnique) (*userproto.ResponseStatus, error) {
	response, err := s.userUsecase.UniqueEmail(ctx, &entity.IsUnique{Email: in.Email})

	if err != nil {
		s.logger.Error(err.Error())
		return &userproto.ResponseStatus{Status: true}, err
	}
	if response.Status {
		return &userproto.ResponseStatus{Status: true}, nil
	}

	return &userproto.ResponseStatus{Status: false}, nil
}

func (s userRPC) UpdateRefresh(ctx context.Context, in *userproto.RefreshRequest) (*userproto.ResponseStatus, error) {
	_, err := s.userUsecase.UpdateRefresh(ctx, &entity.UpdateRefresh{
		UserID:       in.UserId,
		RefreshToken: in.RefreshToken,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return &userproto.ResponseStatus{Status: false}, err
	}

	return &userproto.ResponseStatus{Status: true}, nil
}

func (s userRPC) UpdatePassword(ctx context.Context, in *userproto.UpdatePasswordRequest) (*userproto.ResponseStatus, error) {
	_, err := s.userUsecase.UpdatePassword(ctx, &entity.UpdatePassword{
		UserID:      in.UserId,
		NewPassword: in.NewPassword,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return &userproto.ResponseStatus{Status: false}, err
	}

	return &userproto.ResponseStatus{Status: true}, nil
}
