package mappers

import (
	"user-api/dto"
	"user-api/models"
)

func RegisterReqToUser(req dto.RegisterUserReq) models.User {
	return *models.NewUser(req.Name, req.Age, req.Email, req.Password, req.Address)
}

func UserToPagRes(users []models.User) []dto.UserResponse {
	r := make([]dto.UserResponse, 0)
	for _, u := range users {
		r = append(r, dto.UserResponse{
			Name:  u.Name,
			ID:    u.ID.Hex(),
			Age:   u.Age,
			Email: u.Email,
		})
	}

	return r
}

func UserToRes(user models.User) dto.UserResponse {
	return dto.UserResponse{
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
		ID:    user.ID.Hex(),
	}
}

func UserUpdateReqToUser(user dto.UserUpdateReq) models.User {
	return models.User{
		Name:    user.Name,
		Address: user.Address,
		Age:     user.Age,
	}
}
