package user

type user struct {
	id    string
	name  string
	email string
	token string
}

type InterfaceUser interface {
	GetId() string
	GetName() string
	GetEmail() string
	GetToken() string
}

func NewUser(id string, name string, email string, token string) InterfaceUser {
	return &user{
		id:    id,
		name:  name,
		email: email,
		token: token,
	}
}

func (u *user) GetId() string {
	return u.id
}

func (u *user) GetName() string {
	return u.name
}

func (u *user) GetEmail() string {
	return u.email
}

func (u *user) GetToken() string {
	return u.token
}
