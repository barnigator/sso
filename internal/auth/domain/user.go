package domain

type UserRole string

const (
	UserRoleUser   UserRole = "user"
	UserRoleSeller UserRole = "seller"
	UserRoleAdmin  UserRole = "admin"
)

type User struct {
	ID       int64
	Email    string
	PassHash []byte
	Role     UserRole
	IsActive bool
}
