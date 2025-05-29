package userservice

type UserService interface {
	GetUserByUsername(username string) (*User, error)
	GetUserById(id string) (*User, error)
	UpdateUser(id string, email string, admin bool, displayName string) error
	CreateUser(username string, email string, password string, kind string, admin bool, displayName string) (*User, error)
	ListUsers() ([]User, error)
	Authenticate(username, password string) (*User, error)
	LogoutFromToken(token string) error
	GenerateSessionToken(user *User) (string, error)
	VerifySessionToken(token string) (*Session, *User, error)
}
