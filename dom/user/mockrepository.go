package user

import (
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock

	BaseRepository
}

func (repo *MockRepository) GetUser(userId string) (user *User, err error) {
	args := repo.Called(userId)

	if args.Get(0) != nil {
		user = args.Get(0).(*User)
	}

	return user, args.Error(1)
}

func (repo *MockRepository) PutUser(user *User) error {
	args := repo.Called(user)
	return args.Error(0)
}

func (repo *MockRepository) GetUsersByCountry(cc string) (users []*User, err error) {
	args := repo.Called(cc)

	if args.Get(0) != nil {
		users = args.Get(0).([]*User)
	}

	return users, args.Error(1)
}

func (repo *MockRepository) DeleteUser(id string) error {
	args := repo.Called(id)
	return args.Error(0)
}

func (repo *MockRepository) GetAllUsers() (users []*User, err error) {
	args := repo.Called()

	if args.Get(0) != nil {
		users = args.Get(0).([]*User)
	}

	return users, args.Error(1)
}
