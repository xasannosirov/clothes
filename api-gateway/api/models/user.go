package models

type (
	CreateResponse struct {
		ID string `json:"user_id"`
	}

	UserRegister struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Gender    string `json:"gender"`
	}

	WorkerPost struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Gender      string `json:"gender"`
		Age         int64  `json:"age"`
	}

	WorkerPut struct {
		ID          string `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Gender      string `json:"gender"`
		Age         int64  `json:"age"`
	}

	User struct {
		Id          string `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Gender      string `json:"gender"`
		Age         int64  `json:"age"`
		Refresh     string `json:"refresh_token"`
		Access      string `json:"access_token"`
	}

	ListUser struct {
		User  []User `json:"user"`
		Total uint64 `json:"totcal_count"`
	}

	GoogleUser struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		PictureUrl    string `json:"picture"`
		Locale        string `json:"locale"`
	}
)
