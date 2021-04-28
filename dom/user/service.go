package user

type UserService interface {
	GetUser(string) (*User, error)
	PutUser(*User) (*User, error)
	DeleteUser(string) error
	GetUsersByCountry(cc string) ([]*User, error)
}

type Resources struct {
	Repo Repository
}

type service struct {
	*Resources
}

func NewService(r *Resources) (*service, error) {
	return &service{
		Resources: r,
	}, nil
}
