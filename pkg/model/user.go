package model

type UserStatus int

const (
	UserStatusActive UserStatus = iota
	UserStatusInactive
)

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	Status    UserStatus `json:"status"`
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
}
