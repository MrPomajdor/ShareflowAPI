package auth

import validation "github.com/go-ozzo/ozzo-validation"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m LoginRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, validation.Length(0, 40)),
		validation.Field(&m.Password, validation.Required, validation.Length(0, 128)),
	)
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Authcode  string `json:"authcode"`
}

func (m RegisterRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(0, 40)),
		validation.Field(&m.LastName, validation.Required, validation.Length(0, 40)),
		validation.Field(&m.Email, validation.Required, validation.Length(0, 40)),
		validation.Field(&m.Password, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Authcode, validation.Required, validation.Length(0, 40)),
	)
}
