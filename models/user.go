package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Age      uint8              `bson:"age,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
	Address  string             `bson:"address,omitempty"`
}

func NewUser(name string, age uint8, email string, password string, address string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
		Address:  address,
		Age:      age,
	}
}

func (u *User) HashPassword() (err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return
	}
	u.Password = string(bytes)
	return
}

func (u *User) CheckPassword(providedPassword string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
	if err != nil {
		return
	}
	return
}
