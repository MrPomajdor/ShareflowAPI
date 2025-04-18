package storage

import validation "github.com/go-ozzo/ozzo-validation"

type UploadRequest struct {
	Category string `json:"category_name"`
}

func (m UploadRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Category, validation.Length(0, 128)),
	)
}

type ListRequest struct {
	Category string `json:"category_name"`
}

func (m ListRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Category, validation.Length(0, 128)),
	)
}

type GetRequest struct {
	ID int `json:"id"`
}

func (m GetRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ID, validation.Required, validation.Length(0, 255)),
	)
}

type RemoveRequest struct {
	ID string `json:"id"`
}

func (m RemoveRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ID, validation.Required, validation.Length(0, 255)),
	)
}

type MoveRequest struct {
	ID                  string `json:"id"`
	DestinationCategory string `json:"destination_category"`
}

func (m MoveRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ID, validation.Required, validation.Length(0, 255)),
		validation.Field(&m.DestinationCategory, validation.Required, validation.Length(0, 128)),
	)
}

type CreateURLRequest struct {
	ID string `json:"id"`
}

func (m CreateURLRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ID, validation.Required, validation.Length(0, 255)),
	)
}
