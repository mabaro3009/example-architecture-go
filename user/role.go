package user

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type Role string

func (r Role) String() string {
	return string(r)
}
