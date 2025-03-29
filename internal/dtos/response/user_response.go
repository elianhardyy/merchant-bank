package response

type RegisterResponse struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
	IsActive bool    `json:"is_active"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type RefreshResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Balance  float64  `json:"balance"`
	IsActive bool     `json:"is_active"`
	Roles    []string `json:"role"`
}
