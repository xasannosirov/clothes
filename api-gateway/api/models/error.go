package models

// Error ...
type Error struct {
	Message string `json:"message"`
}

const (
	// Error messages
	RequiredRefreshMessage          = "Required refresh"
	NoAccessMessage                 = "You have no access this page"
	TokenInvalidMessage             = "Invalid token claims"
	IncorrectLoginMessage           = "Incorrect login"
	IncorrectPasswordMessage        = "Incorrect password"
	IncorrectLoginOrPassMessage     = "Incorrect login or password"
	IncorrectConfirmPasswordMessage = "Password not equal"
	WeakPasswordMessage             = "Password is weak"
	IncorrectDateFormatMessage      = "Date was not valid"
	WrongInfoMessage                = "Data was not valid"
	LimitErrorMessgae               = "Exceeded the limit"
	EmailUsedMessage                = "This email already in use"
	InvalidEmailMessage             = "Email invalid"
	InvalidPhoneNumberMessage       = "Phone number invalid"
	InvalidGenderMessage            = "Gender invalid"
	NotMatchOTP                     = "Verification code not match"
	NotMachPassword                 = "Password not match"
	NotFoundMessage                 = "Data was not found"
	NotAddedMessage                 = "Data was not added"
	NotUpdatedMessage               = "Data was not updated"
	NotDeletedMessage               = "Data was not deleted"
	NotCreatedMessage               = "Data was not created"
	InternalMessage                 = "Error happened during process"
	TokenExchangeMessage            = "Token exchange failed"
	FailedGetUserInfo               = "Failed to get user info from google"

	// Google
	MisMatchMessage = "State mismatch"
	MisCode         = "Missing code"

	// Info messages
	SentOTPMessage         = "We have send verification code your email"
	SuccessDeleted         = "Data successfully deleted"
	SuccessAdded           = "Data successfully added"
	SuccessUpdatedPassword = "Password successfully updated"
)
