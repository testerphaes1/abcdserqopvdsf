package usecase_models

type Auth struct {
	Email                 string `json:"email"`
	EmailVerificationCode string `json:"email_verification_code"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AccountAuth struct {
}
