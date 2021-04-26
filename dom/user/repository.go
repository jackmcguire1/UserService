package user

import "fmt"

var NotImplementedErr = fmt.Errorf("this method is not implemented")

type Repository interface {
	GetUser(string) (*User, error)
	GetAllUsers(string, int) ([]*User, string, error)
	PutUser(*User) error
}

type BaseRepository struct{}

func (repo *BaseRepository) GetUser(string) (*User, error) {
	return nil, NotImplementedErr
}

func (repo *BaseRepository) PutUser(*User) error {
	return NotImplementedErr
}

func (repo *BaseRepository) GetAllUsers(cursor string, limit int) ([]*User, string, error) {
	return nil, "", NotImplementedErr
}
