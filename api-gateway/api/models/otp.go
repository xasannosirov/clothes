package models

type Otp struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}