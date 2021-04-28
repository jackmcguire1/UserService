package user

import "fmt"

var NotImplementedErr = fmt.Errorf("this method is not implemented")

type Repository interface {
	GetUser(string) (*User, error)
	GetUsersByCountry(cc string) (users []*User, err error)
	DeleteUser(string) error
	PutUser(*User) error
}

type BaseRepository struct{}

func (repo *BaseRepository) GetUser(string) (*User, error) {
	return nil, NotImplementedErr
}

func (repo *BaseRepository) PutUser(*User) error {
	return NotImplementedErr
}

func (repo *BaseRepository) DeleteUser(string) error {
	return NotImplementedErr
}

func (repo *BaseRepository) GetUsersByCountry(cc string) (users []*User, err error) {
	return nil, NotImplementedErr
}
