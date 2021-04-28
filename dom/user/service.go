package user

type UserUpdate struct {
	User   *User
	Status string
}

type UserService interface {
	GetUser(string) (*User, error)
	PutUser(*User) (*User, error)
	DeleteUser(string) error
	GetUsersByCountry(string) ([]*User, error)
}

type Resources struct {
	Repo        Repository
	UserChannel chan *UserUpdate
}

type service struct {
	*Resources
}

func NewService(r *Resources) (*service, error) {
	return &service{
		Resources: r,
	}, nil
}
