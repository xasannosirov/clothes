package models

type CreateUser struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponseModel struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type JwtRequestModel struct {
	Token string `json:"token"`
}

type ResponseError struct {
	Error interface{} `json:"error"`
}

// ServerError ...
type ServerError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserCreateReq struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}
type UserResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Password     string `json:"password"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
type UserList struct {
	ListUser []UserResponse
}
