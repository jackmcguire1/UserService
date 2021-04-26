package user

type UserService interface {
	GetUser(string) (*User, error)
	PutUser(*User) (*User, error)
	GetUsersByCountry(cc string, cursor string, limit int) ([]*User, string, error)
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
