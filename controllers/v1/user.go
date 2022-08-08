package controllers

import (
	"log"
	"net/http"
	"strconv"
	"user-api/dto"
	"user-api/mappers"
	"user-api/models"
	"user-api/response"
	"user-api/services"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetAll() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Update() gin.HandlerFunc
	GetById() gin.HandlerFunc
}

type UserControllerImpl struct {
	svc services.UserService
}

func NewUserJson(svc services.UserService) UserController {
	return UserControllerImpl{svc: svc}
}

// Get example godoc
// @SummaryUser Get Users paginated
// @Description Get all users paginated
// @Param limit query integer false "limit"
// @Param page query integer false "page"
// @Param token header string true "Authentication token"
// @Accept json
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Router /users [get]
func (u UserControllerImpl) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := strconv.ParseUint(c.Query("limit"), 0, 64)

		if err != nil || limit == 0 || limit > 10 {
			limit = 10
		}

		page, err := strconv.ParseUint(c.Query("page"), 0, 64)

		if err != nil || page == 0 {
			page = 1
		}

		u, apiErr := u.svc.GetAll(limit, page)

		if apiErr.Status != 0 {
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}
		resp := mappers.UserToPagRes(u)
		c.AbortWithStatusJSON(http.StatusOK, response.ApiResponse{Body: resp, Status: http.StatusOK})
	}
}

// Delete example godoc
// @SummaryUser Delete user
// @Description Delete user by id
// @Param id path string true "id"
// @Produce json
// @Success 204
// @Router /users/:id [delete]
func (u UserControllerImpl) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, exists := c.Get("user")

		if !exists {
			log.Printf("[USER CONTROLLER] User not found in context with email")
			e := response.InternalServerError
			c.AbortWithStatusJSON(e.Status, e)
			return
		}

		if user.(models.User).ID.Hex() != id {
			log.Println("[USER CONTROLLER] Cannot delete different user")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err := u.svc.DeleteById(id)

		if err.Status != 0 {
			c.AbortWithStatusJSON(err.Status, err)
			return
		}

		c.AbortWithStatus(http.StatusNoContent)
	}
}

// Update example godoc
// @SummaryUser Update user
// @Description Update user by id
// @Param Update body dto.UserUpdateReq true "Update request"
// @Produce json
// @Success 200
// @Router /users/:id [put]
func (u UserControllerImpl) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.UserUpdateReq{}
		id := c.Param("id")

		err := c.ShouldBindJSON(&req)

		if err != nil {
			log.Printf("Error parsing user input error: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, response.BadRequestError)
			return
		}

		user, exists := c.Get("user")

		if !exists {
			log.Printf("[USER CONTROLLER] User not found in context with email")
			e := response.InternalServerError
			c.AbortWithStatusJSON(e.Status, e)
			return
		}

		if user.(models.User).ID.Hex() != id {
			log.Println("[USER CONTROLLER] Cannot update different user")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		v := req.ValidateFields()

		if len(v) != 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, v)
			return
		}

		apiErr := u.svc.UpdateById(id, mappers.UserUpdateReqToUser(req))

		if apiErr.Status != 0 {
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
		}

		c.Status(http.StatusNoContent)
	}
}

// Get user example godoc
// @SummaryUser Get User by id
// @Description Get user by id
// @Param Id path string true "User id"
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Router /users/:id [get]
func (u UserControllerImpl) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		u, err := u.svc.FindById(id)

		if err.Status != 0 {
			ctx.AbortWithStatusJSON(err.Status, err)
			return
		}

		uRes := mappers.UserToRes(u)
		ctx.JSON(http.StatusOK, uRes)
	}
}
