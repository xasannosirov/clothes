package postgresql

// import (
// 	"context"
// 	"github.com/google/uuid"
// 	"log"
// 	"testing"
// 	pbu "user-service/genproto/user_service"
// 	"user-service/internal/entity"
// 	"user-service/internal/infrastructure/repository"
// 	"user-service/internal/pkg/config"
// 	db "user-service/internal/pkg/postgres"

// 	"github.com/stretchr/testify/suite"
// )

// type UserRepositorySuiteTest struct {
// 	suite.Suite
// 	Repository repository.Users
// }

// func (p *UserRepositorySuiteTest) SetupSuite() {
// 	pgPoll, err := db.New(config.New())
// 	if err != nil {
// 		log.Fatal("Error while connecting database with suite test!")
// 		return
// 	}

// 	p.Repository = NewUsersRepo(pgPoll)
// }

// func (u *UserRepositorySuiteTest) TestSuite() {

// 	// mock info for check user-service methods
// 	defaultGUID := uuid.New().String()
// 	userReq := &pbu.User{
// 		Id:          defaultGUID,
// 		FirstName:   "Davron",
// 		LastName:    "Nuriddinov",
// 		Email:       "nuriddinovdavron2003@gmail.com",
// 		Password:    "Secret12345",
// 		PhoneNumber: "+998557779944",
// 	}

// 	// Create user method
// 	createdUser, err := u.Repository.Create(context.TODO(), &entity.User{
// 		GUID:        userReq.Id,
// 		FirstName:   userReq.FirstName,
// 		LastName:    userReq.LastName,
// 		Email:       userReq.Email,
// 		Password:    userReq.Password,
// 		PhoneNumber: userReq.PhoneNumber,
// 	})

// 	u.Suite.NoError(err)
// 	u.Suite.NotNil(createdUser)
// 	u.Suite.NotNil(createdUser.GUID)
// 	u.Suite.Equal(createdUser.GUID, userReq.Id)
// 	u.Suite.Equal(createdUser.FirstName, userReq.FirstName)
// 	u.Suite.Equal(createdUser.LastName, userReq.LastName)
// 	u.Suite.Equal(createdUser.Email, userReq.Email)
// 	u.Suite.Equal(createdUser.Password, userReq.Password)
// 	u.Suite.Equal(createdUser.PhoneNumber, userReq.PhoneNumber)

// 	// update user method
// 	updatedUser, err := u.Repository.Update(context.Background(), &entity.User{
// 		GUID:        createdUser.GUID,
// 		FirstName:   "New " + userReq.FirstName,
// 		LastName:    "New " + userReq.LastName,
// 		Email:       "New " + userReq.Email,
// 		Password:    "New " + userReq.Password,
// 		PhoneNumber: "New " + userReq.PhoneNumber,
// 	})
// 	u.Suite.NoError(err)
// 	u.Suite.NotNil(updatedUser)
// 	u.Suite.NotEmpty(updatedUser)
// 	u.Suite.Equal(updatedUser.GUID, createdUser.GUID)
// 	u.Suite.Equal(updatedUser.FirstName, "New "+userReq.FirstName)
// 	u.Suite.Equal(updatedUser.LastName, "New "+userReq.LastName)
// 	u.Suite.Equal(updatedUser.Email, "New "+userReq.Email)
// 	u.Suite.Equal(updatedUser.Password, "New "+userReq.Password)
// 	u.Suite.Equal(updatedUser.PhoneNumber, "New "+userReq.PhoneNumber)

// 	// get user method
// 	getUser, err := u.Repository.Get(context.TODO(), map[string]string{
// 		"id": updatedUser.GUID,
// 	})
// 	u.Suite.NoError(err)
// 	u.Suite.Equal(getUser.GUID, updatedUser.GUID)
// 	u.Suite.Equal(getUser.FirstName, updatedUser.FirstName)
// 	u.Suite.Equal(getUser.LastName, updatedUser.LastName)
// 	u.Suite.Equal(getUser.Email, updatedUser.Email)
// 	u.Suite.Equal(getUser.Password, updatedUser.Password)
// 	u.Suite.Equal(getUser.PhoneNumber, updatedUser.PhoneNumber)

// 	// list user method
// 	listUserResponse, err := u.Repository.List(context.TODO(), 10, 1, map[string]string{})
// 	u.Suite.NoError(err)
// 	u.Suite.Equal(listUserResponse[0].FirstName, getUser.FirstName)
// 	u.Suite.Equal(listUserResponse[0].LastName, getUser.LastName)
// 	u.Suite.Equal(listUserResponse[0].Email, getUser.Email)
// 	u.Suite.Equal(listUserResponse[0].Email, getUser.Email)
// 	u.Suite.Equal(listUserResponse[0].Password, getUser.Password)
// 	u.Suite.Equal(listUserResponse[0].PhoneNumber, getUser.PhoneNumber)

// 	// Unique Email
// 	errResponse, err := u.Repository.UniqueEmail(context.TODO(), &entity.IsUnique{
// 		Email: updatedUser.Email,
// 	})
// 	u.Suite.Error(err)
// 	u.Suite.Equal(errResponse.Status, false)
// 	successResponse, err := u.Repository.UniqueEmail(context.TODO(), &entity.IsUnique{
// 		Email: "xasannosriov094@gmail.com",
// 	})
// 	u.Suite.NoError(err)
// 	u.Suite.Equal(successResponse.Status, true)

// 	// Update Refresh Token
// 	refreshStatus, err := u.Repository.UpdateRefresh(context.TODO(), &entity.UpdateRefresh{
// 		UserID:       listUserResponse[0].GUID,
// 		Role:         "user",
// 		RefreshToken: "new.secret.refreshtoken",
// 	})
// 	u.Suite.NoError(err)
// 	u.Suite.Equal(refreshStatus.Status, true)

// 	// Update Password
// 	passwordStatus, err := u.Repository.UpdatePassword(context.TODO(), &entity.UpdatePassword{
// 		UserID:      listUserResponse[0].GUID,
// 		Role:        "user",
// 		NewPassword: "SecretPassword",
// 	})
// 	u.Suite.NoError(err)
// 	u.Suite.Equal(passwordStatus.Status, true)

// 	// delete user method
// 	err = u.Repository.Delete(context.TODO(), getUser.GUID)
// 	u.Suite.NoError(err)
// }

// // running user-service suite test
// func TestUserRepositoryTestSuite(t *testing.T) {
// 	suite.Run(t, new(UserRepositorySuiteTest))
// }
