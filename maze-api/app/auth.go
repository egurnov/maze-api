package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	CTXUserID = "user_id"
)

type LoginCredentialsDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary Login to the API
// @Description Login with username and password, get a token to use for further operations.
// @ID login
// @Tags Auth
// @Accept json
// @Produce json
// @Param email body Credentials true "Credentials"
// @Success 200 {object} JWTTokenResp
// @Failure 400 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /login [post]
func (a *App) Login(ctx *gin.Context) {
	var credentials LoginCredentialsDTO
	err := ctx.ShouldBind(&credentials)
	if err != nil {
		ctx.Error(err).SetType(BadRequestErrorType)
		return
	}

	user, err := a.UserService.GetByUsername(credentials.Username)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	token, err := a.JWTService.GenerateToken(user.ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, LoginResponseDTO{
		Token: token,
	})
}

func (a *App) AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := a.JWTService.ValidateToken(tokenString)
		if err != nil {
			a.Log.Debug("Failed authentication attempt: ", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set(CTXUserID, claims.UserID)
		a.Log.WithFields(map[string]interface{}{
			"UserID": claims.UserID,
		}).Debugf("Successful authentication")
	}
}
