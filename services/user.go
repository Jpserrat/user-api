package services

import (
	"log"
	"user-api/auth"
	"user-api/models"
	"user-api/repositories"
	"user-api/response"
)

type UserService interface {
	Register(models.User) (models.User, response.ApiError)
	GetAll(uint64, uint64) ([]models.User, response.ApiError)
	FindByEmail(email string) (models.User, response.ApiError)
	FindById(id string) (models.User, response.ApiError)
	DeleteById(id string) response.ApiError
	UpdateById(id string, u models.User) response.ApiError
	Login(email string, password string) (string, response.ApiError)
}

type userServiceImpl struct {
	r repositories.UserRepo
}

func NewUser(r repositories.UserRepo) UserService {
	return userServiceImpl{
		r: r,
	}
}

func (svc userServiceImpl) Register(u models.User) (models.User, response.ApiError) {

	_, apiErr := svc.FindByEmail(u.Email)

	if apiErr.Status != 0 {
		if apiErr.Status != response.ResourceNotFoundError.Status {
			log.Printf("[USER SERVICE] Unexepcted error ocurred: %s", apiErr.Error)
			return u, response.InternalServerError
		}
	} else {
		log.Printf("[USER SERVICE] Email %s already in use", u.Email)
		return u, response.EmailAlreadyInUse
	}

	err := u.HashPassword()
	if err != nil {
		log.Printf("[USER SERVICE] Error hashing password: %s", err.Error())
		return u, response.InternalServerError
	}

	return u, svc.r.Save(u)
}

func (svc userServiceImpl) GetAll(limit uint64, page uint64) ([]models.User, response.ApiError) {
	u, err := svc.r.GetAll(limit, page)
	if err != nil {
		log.Printf("Error getting users from repository: %v", err.Error())
		return u, response.InternalServerError
	}

	return u, response.ApiError{}
}

func (svc userServiceImpl) Login(email, password string) (string, response.ApiError) {
	u, apiErr := svc.FindByEmail(email)

	if apiErr.Status != 0 {
		if apiErr.Status == response.ResourceNotFoundError.Status {
			log.Printf("[USER SERVICE] User not found with email: %s", email)
		}
		return "", apiErr
	}

	err := u.CheckPassword(password)

	if err != nil {
		log.Printf("[USER SERVICE] Invalid password")
		return "", response.InvalidCredentialsError
	}

	jwt, err := auth.GenerateJWT(email)

	if err != nil {
		log.Printf("[USER SERVICE] Error generating JWT: %s", err.Error())
		return "", response.InternalServerError
	}

	return jwt, response.ApiError{}
}

func (svc userServiceImpl) FindByEmail(email string) (models.User, response.ApiError) {
	return svc.r.FindByField(email, "email")
}

func (svc userServiceImpl) FindById(id string) (models.User, response.ApiError) {

	return svc.r.FindById(id)
}

func (svc userServiceImpl) DeleteById(id string) response.ApiError {

	return svc.r.DeleteById(id)
}

func (svc userServiceImpl) UpdateById(id string, u models.User) response.ApiError {

	return svc.r.UpdateByID(id, u)
}
