package models

type User struct {
	Id          string
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber string
	Gender      string
	Age         int64
	Role        string
	Refresh     string
}
type UserCreateReq struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber string
	Gender      string
	Age         int64
}
type UserResponse struct {
	Id          string
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber string
	Gender      string
	Age         int64
	Role        string
	Refresh     string
	Access      string
}
type Response struct{
	Response string
}
type VerifyForgetPassword struct{
	Otp string
	Email string
	NewPassword string
}