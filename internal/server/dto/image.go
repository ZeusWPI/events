package dto

type ImageSave struct {
	Name string `form:"name" validate:"required"`
	File []byte `validate:"required,min=1"`
}
