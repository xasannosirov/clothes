package models

type (
	UserRegister struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Gender    string `json:"gender"`
		Role      string `json:"role"`
	}

	User struct {
		Id          string `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Gender      string `json:"gender"`
		Age         int64  `json:"age"`
		Role        string `json:"role"`
		Refresh     string `json:"refresh_token"`
		Access      string `json:"access_token"`
	}
)
