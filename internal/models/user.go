package models

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type User struct {
	ID           int    `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	Who          string `json:"who" db:"who"`
	Name         string `json:"name" db:"name"`
	PassHash     string `json:"pass_hash" db:"pass_hash"`
	IsRegistered bool   `json:"is_registered" db:"is_registered"`
}

type SignUp struct {
	Name     string `json:"name" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (i SignUp) Validate() error {
	return validate.Struct(i)
}

type SignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (i SignIn) Validate() error {
	return validate.Struct(i)
}
