package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/egurnov/maze-api/maze-api/model"
)

type CreateUserRequestDTO struct {
	Username string `json:"username" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty,min=1"`
}

type IDResponseDTO struct {
	ID int64 `json:"id"`
}

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
// @Param credentials body CreateUserRequestDTO true "New user credentials"
// @Success 201 {object} IDResponseDTO
// @Failure 400 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /user [post]
func (a *App) CreateUser(ctx *gin.Context) {
	var user CreateUserRequestDTO
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

	ctx.JSON(http.StatusCreated, IDResponseDTO{ID: id})
}
