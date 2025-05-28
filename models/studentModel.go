package models


type Student struct {
	Id              int64  `json:"id"`
	FullName        string `json:"full_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password,omitempty" validate:"required,min=6"` // omitted in JSON response
	Age             int    `json:"age" validate:"required,min=16,max=100"`
	Gender          string `json:"gender" validate:"required,oneof=male female other"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	DOB             string `json:"dob" validate:"required"` // format: YYYY-MM-DD
	Address         string `json:"address" validate:"required"`
}
