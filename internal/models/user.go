package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
