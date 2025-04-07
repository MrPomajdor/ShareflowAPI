package storage

import validation "github.com/go-ozzo/ozzo-validation"

type UploadRequest struct {
	Path string `json:"path"`
}

func (m UploadRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Path, validation.Required, validation.Length(0, 255)),
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
	Path        string `json:"path"`
	Destination string `json:"destination"`
}

func (m MoveRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Path, validation.Required, validation.Length(0, 255)),
		validation.Field(&m.Destination, validation.Required, validation.Length(0, 255)),
	)
}

type CreateURLRequest struct {
	Path string `json:"path"`
}

func (m CreateURLRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Path, validation.Required, validation.Length(0, 255)),
	)
}
