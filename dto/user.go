package dto

import (
	"net/url"

	"github.com/thedevsaddam/govalidator"
)

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   uint8  `json:"age"`
}

type UserUpdateReq struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Age     uint8  `json:"age" validate:"required"`
}

func (req UserUpdateReq) ValidateFields() url.Values {
	rules := govalidator.MapData{
		"name":    []string{"required", "min:3"},
		"age":     []string{"required"},
		"address": []string{"required"},
	}

	opts := govalidator.Options{
		Data:  &req,
		Rules: rules,
	}
	v := govalidator.New(opts)

	return v.ValidateStruct()
}
