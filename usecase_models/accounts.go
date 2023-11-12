package usecase_models

type RegisterAccountRequest struct {
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	PhoneNumber           string `json:"phone_number"`
	Email                 string `json:"email"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	EmailVerificationCode string `json:"email_verification_code"`
}
type Account struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type AccountResponse struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Username    string `json:"username"`
}

type RegisterAccountResponse struct {
	AccountId int    `json:"account_id"`
	Token     string `json:"token"`
}

type EmailVerificationRequest struct {
	Email string `json:"email"`
}
