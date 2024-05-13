package models

type (
	UserCreateResponse struct {
		UserID string `json:"user_id"`
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
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Gender      string `json:"gender"`
		Age         int64  `json:"age"`
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
		Refresh     string `json:"refresh_token"`
		Access      string `json:"access_token"`
	}

	ListUser struct {
		User  []User `json:"user"`
		Total uint64 `json:"totcal_count"`
	}
)
