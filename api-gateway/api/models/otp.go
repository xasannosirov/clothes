package models

type (
	Login struct {
		Email    string `json:"email" example:"xasannosirov094@gmail.com"`
		Password string `json:"password" example:"Sehtols@01"`
	}

	Otp struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	ResetPassword struct {
		Otp         string `json:"otp"`
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}

	LoginResp struct {
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

	TokenResp struct {
		ID      string `json:"user_id"`
		Access  string `json:"access_token"`
		Refresh string `json:"refresh_token"`
		Role    string `json:"role"`
	}
)
