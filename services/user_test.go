package services

import (
	"testing"
	mocks "user-api/mocks/repositories"
	"user-api/models"
	"user-api/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterAlreadyExists(t *testing.T) {
	userToBeRegister := models.NewUser("test", 20, "test@test.com", "pass", "add")
	mockUserRepo := new(mocks.UserRepo)
	svc := userServiceImpl{r: mockUserRepo}
	mockUserRepo.On("FindByField", userToBeRegister.Email, "email").Return(*userToBeRegister, response.ApiError{})

	_, apiErr := svc.Register(*userToBeRegister)

	assert.Equal(t, response.EmailAlreadyInUse.Code, apiErr.Code)
}

func TestRegisterInternalError(t *testing.T) {
	userToBeRegister := models.NewUser("test", 20, "test@test.com", "pass", "add")
	mockUserRepo := new(mocks.UserRepo)
	svc := userServiceImpl{r: mockUserRepo}
	mockUserRepo.On("FindByField", userToBeRegister.Email, "email").Return(*userToBeRegister, response.InternalServerError)

	_, apiErr := svc.Register(*userToBeRegister)

	assert.Equal(t, response.InternalServerError.Code, apiErr.Code)
}

func TestRegisterSuccess(t *testing.T) {
	userToBeRegister := models.NewUser("test", 20, "test@test.com", "pass", "add")
	mockUserRepo := new(mocks.UserRepo)
	svc := userServiceImpl{r: mockUserRepo}
	mockUserRepo.On("FindByField", userToBeRegister.Email, "email").Return(*userToBeRegister, response.ResourceNotFoundError)
	mockUserRepo.On("Save", mock.AnythingOfType("models.User")).Return(response.ApiError{})

	user, apiErr := svc.Register(*userToBeRegister)

	assert.Nil(t, user.CheckPassword("pass"))
	assert.Equal(t, "", apiErr.Code)
}

func TestLoginUserNotFound(t *testing.T) {
	email := "test@test.com"
	password := "test"
	mockUserRepo := new(mocks.UserRepo)
	svc := userServiceImpl{r: mockUserRepo}
	mockUserRepo.On("FindByField", email, "email").Return(models.User{}, response.ResourceNotFoundError)

	_, err := svc.Login(email, password)

	assert.Equal(t, response.ResourceNotFoundError.Code, err.Code)
}

func TestLoginInvalidCredentials(t *testing.T) {
	email := "test@test.com"
	password := "test"
	user := models.User{Password: "tes"}
	user.HashPassword()
	mockUserRepo := new(mocks.UserRepo)
	svc := userServiceImpl{r: mockUserRepo}
	mockUserRepo.On("FindByField", email, "email").Return(user, response.ApiError{})

	_, err := svc.Login(email, password)

	assert.Equal(t, response.InvalidCredentialsError.Code, err.Code)
}

func TestLoginSuccess(t *testing.T) {
	email := "test@test.com"
	password := "test"
	user := models.User{Password: "test"}
	user.HashPassword()
	mockUserRepo := new(mocks.UserRepo)
	svc := userServiceImpl{r: mockUserRepo}
	mockUserRepo.On("FindByField", email, "email").Return(user, response.ApiError{})

	jwt, err := svc.Login(email, password)

	assert.Equal(t, 0, err.Status)
	assert.NotEmpty(t, jwt)
}

func TestShouldCallFindById(t *testing.T) {
	id := "id"
	mockUserRepo := new(mocks.UserRepo)
	svc := userServiceImpl{r: mockUserRepo}
	user := models.User{Name: "test"}
	mockUserRepo.On("FindById", id).Return(user, response.ApiError{})

	u, _ := svc.FindById(id)

	assert.Equal(t, user.Email, u.Email)
}

func TestShouldCallDeleteById(t *testing.T) {
	id := "id"
	mockUserRepo := new(mocks.UserRepo)
	svc := userServiceImpl{r: mockUserRepo}
	err := response.ApiError{Code: "CODE"}
	mockUserRepo.On("DeleteById", id).Return(err)

	apiErr := svc.DeleteById(id)

	assert.Equal(t, err.Code, apiErr.Code)
}

func TestShouldCallUpdateById(t *testing.T) {
	id := "id"
	mockUserRepo := new(mocks.UserRepo)
	svc := userServiceImpl{r: mockUserRepo}
	user := models.User{Name: "test"}
	err := response.ApiError{Code: "CODE"}
	mockUserRepo.On("UpdateByID", id, user).Return(err)

	apiErr := svc.UpdateById(id, user)

	assert.Equal(t, err.Code, apiErr.Code)
}
