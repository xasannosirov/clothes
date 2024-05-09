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
		ID      string `json:"user_id"`
		Role    string `json:"role"`
		Acccess string `json:"access_token"`
		Refresh string `json:"refresh_token"`
		Gender  string `json:"gender"`
		Age     string `json:"age"`
	}

	TokenResp struct {
		ID      string `json:"user_id"`
		Access  string `json:"access_token"`
		Refresh string `json:"refresh_token"`
		Role    string `json:"role"`
	}
)
