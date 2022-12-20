package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/egurnov/maze-api/maze-api/model"
)

func (a *App) GetUserByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 0, 64)
	if err != nil {
		ctx.Error(err).SetType(BadRequestErrorType)
		return
	}

	res, err := a.UserService.GetByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// CreateUser godoc
// @Summary Register a new user for the API
// @Description Provide a unique username and password to create a new user.
// @ID CreateUser
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body Credentials true "New user credentials"
// @Success 201 {object} IDResponse
// @Failure 400 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /post [post]
func (a *App) CreateUser(ctx *gin.Context) {
	var user CreateUserDTO
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.Error(err).SetType(BadRequestErrorType)
		return
	}

	if len(user.Password) == 0 {
		ctx.Error(model.ErrInvalidInput)
		return
	}

	id, err := a.UserService.Create(&model.User{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, IDResponse{ID: id})
}

type CreateUserDTO struct {
	Username string `json:"username,omitempty" binding:"omitempty"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}
