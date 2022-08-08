package dto

import (
	"net/url"

	"github.com/thedevsaddam/govalidator"
)

type RegisterUserReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Age      uint8  `json:"age" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Address  string `json:"address" validate:"required"`
}

func (req RegisterUserReq) ValidateFields() url.Values {
	rules := govalidator.MapData{
		"name":     []string{"required", "min:3"},
		"email":    []string{"required", "min:4", "email"},
		"age":      []string{"required"},
		"password": []string{"required", "min:6"},
		"address":  []string{"required"},
	}

	opts := govalidator.Options{
		Data:  &req,
		Rules: rules,
	}
	v := govalidator.New(opts)

	return v.ValidateStruct()
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (req LoginReq) ValidateFields() url.Values {
	rules := govalidator.MapData{
		"email":    []string{"required", "min:4", "email"},
		"password": []string{"required", "min:6"},
	}

	opts := govalidator.Options{
		Data:  &req,
		Rules: rules,
	}
	v := govalidator.New(opts)

	return v.ValidateStruct()
}

type LoginRes struct {
	Jwt string `json:"jwt"`
}
