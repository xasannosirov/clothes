package postgresql

import (
	"context"
	"log"
	pbu "user-service/genproto/user_service"
	"user-service/internal/entity"
	"user-service/internal/infrastructure/repository"
	"user-service/internal/pkg/config"
	db "user-service/internal/pkg/postgres"

	"github.com/stretchr/testify/suite"
)

type UserRepositorySuiteTest struct {
	suite.Suite
	Repository repository.Users
}

func (p *UserRepositorySuiteTest) SetupSuite(){
	pgPoll, err := db.New(config.New())
	if err != nil{
		log.Fatal("Error while connecting database with suite test!")
	}

	p.Repository = NewUsersRepo(pgPoll)
}

func (u *UserRepositorySuiteTest) TestUserCRUD(){
	userReq := &pbu.User{
		FirstName: "",
		LastName: "test",
		Email: "nuriddinovdavron2003@gmail.com",
		Password: "test",
		PhoneNumber: "test",
		Gender: "test",
		Age: 1,
	}

	//Create method
	err := u.Repository.Create(context.Background(), &entity.User{
		FirstName: userReq.FirstName,
	})
	u.Suite.NoError(err)

	


}
