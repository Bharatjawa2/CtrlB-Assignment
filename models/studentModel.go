package models


type Student struct {
	ID        int    `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `validate:"required"`
	Password  string `validate:"required"`
	Phone     string `validate:"required"`
	Address   string `validate:"required"`
}
