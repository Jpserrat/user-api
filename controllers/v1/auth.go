package controllers

import (
	"log"
	"net/http"
	"user-api/auth"
	"user-api/dto"
	"user-api/mappers"
	"user-api/response"
	"user-api/services"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	VerifyToken() gin.HandlerFunc
}

type AuthControllerImpl struct {
	userSvc services.UserService
}

func NewAuth(uSvc services.UserService) AuthController {
	return AuthControllerImpl{
		userSvc: uSvc,
	}
}

func (a AuthControllerImpl) VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			log.Println("[AUTH CONTROLLER] Empty token")
			e := response.InvalidTokenError
			ctx.AbortWithStatusJSON(e.Status, e)
			return
		}

		email, err := auth.ValidateToken(token)
		if err != nil {
			log.Printf("[AUTH CONTROLLER] Invalid token %s", token)
			e := response.InvalidTokenError
			ctx.AbortWithStatusJSON(e.Status, e)
			return
		}

		user, apiErr := a.userSvc.FindByEmail(email)

		if apiErr.Status != 0 {
			ctx.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}
		ctx.Set("user", user)

		ctx.Next()
	}
}

// Login example godoc
// @SummaryUser login
// @Description do login
// @Param Login body dto.LoginReq true "User credentials"
// @Accept json
// @Produce json
// @Success 200 {object} dto.LoginRes
// @Router /auth/login [post]
func (a AuthControllerImpl) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		req := dto.LoginReq{}
		err := ctx.ShouldBindJSON(&req)

		if err != nil {
			log.Printf("Error parsing user input error: %s", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.BadRequestError)
			return
		}

		v := req.ValidateFields()

		if len(v) != 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, v)
			return
		}

		jwt, apiErr := a.userSvc.Login(req.Email, req.Password)

		if apiErr.Status != 0 {
			ctx.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}

		ctx.JSON(http.StatusOK, dto.LoginRes{Jwt: jwt})
	}
}

// Register example godoc
// @SummaryUser Register
// @Description Register a new user
// @Param Register body dto.RegisterUserReq true "User information"
// @Accept json
// @Produce json
// @Success 201
// @Router /auth/register [post]
func (a AuthControllerImpl) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := dto.RegisterUserReq{}
		err := ctx.ShouldBindJSON(&req)

		if err != nil {
			log.Printf("Error parsing user input error: %s", err.Error())
			ctx.JSON(http.StatusBadRequest, response.BadRequestError)
			ctx.Abort()
			return
		}

		v := req.ValidateFields()

		if len(v) != 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, v)
			return
		}

		user := mappers.RegisterReqToUser(req)

		log.Printf("Register request mapped to user %v", user)

		_, apiErr := a.userSvc.Register(user)

		if apiErr.Status != 0 {
			log.Printf("Error register user: %v", req)
			ctx.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}

		ctx.AbortWithStatus(http.StatusCreated)
	}

}
