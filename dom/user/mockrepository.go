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

func (repo *MockRepository) GetAllUsers(cursor string, limit int) (users []*User, bookmark string, err error) {
	args := repo.Called(cursor, limit)

	if args.Get(0) != nil {
		users = args.Get(0).([]*User)
	}

	return users, args.String(1), args.Error(2)
}
