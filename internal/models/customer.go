package models

type Customer struct {
	ID       string  `json:"id"`
	Username string  `json:"username" validate:"required,min=5,alphanum,username_check"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8,password_check"`
	Balance  float64 `json:"balance"`
	IsActive bool    `json:"is_active"`
}
