package info

import validation "github.com/go-ozzo/ozzo-validation"

type UpdateDataRequest struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (m UpdateDataRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Field, validation.Required, validation.Length(0, 40)),
		validation.Field(&m.Value, validation.Required, validation.Length(0, 40)),
	)
}
