package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/egurnov/maze-api/maze-api/model"
)

const (
	BadRequestErrorType gin.ErrorType = 3
	NotFoundErrorType   gin.ErrorType = 4
)

type Message struct {
	Message string `json:"message"`
}

var MessageOK = &Message{
	Message: "OK",
}

func (a *App) renderErrors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, err := range ctx.Errors {
			switch err.Err {
			case nil:
				continue
			case model.ErrNotFound:
				ctx.JSON(http.StatusNotFound, &Message{Message: err.Error()})
			case model.ErrInvalidInput, model.ErrUsernameAlreadyUsed:
				ctx.JSON(http.StatusBadRequest, &Message{Message: err.Error()})
			case model.ErrNotAllowed:
				ctx.JSON(http.StatusForbidden, &Message{Message: err.Error()})
			case model.ErrorUnauthorized:
				ctx.JSON(http.StatusUnauthorized, &Message{Message: err.Error()})
			default:
				switch err.Type {
				case BadRequestErrorType, NotFoundErrorType:
					ctx.JSON(http.StatusBadRequest, &Message{Message: err.Error()})
				default:
					a.Log.WithField("url", ctx.Request.URL.String()).Error(err)
					ctx.JSON(http.StatusInternalServerError, &Message{
						Message: "Internal Server Error",
					})
				}
			}
		}
	}
}
