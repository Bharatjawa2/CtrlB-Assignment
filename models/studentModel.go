package models


type Student struct {
	Id    int64
	Name  string `validate:"required"`
	Email string `validate:"required"`
	AGE   int    `validate:"required"`
}
