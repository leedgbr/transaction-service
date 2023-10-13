package errorhandling

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"transaction-service/internal/business"
)

const (
	BadRequest = "BAD_REQUEST"

	systemErrorMessage     string = "SYSTEM_ERROR"
	badRequestErrorMessage string = "BAD_REQUEST"
)

// ErrorResponse represents the http response body for an error.
type ErrorResponse struct {
	Message string `json:"message"`
}

// NewMiddleware is middleware for gin that provides top level error handling.  It is responsible for making sure the
// http response and status are appropriate for the error(s) that occurred.  Principally it distinguishes between
// business and system errors, with business errors resulting in a 422 http status and system errors resulting in a 500
// http status.  System errors return a static error message, with details logged on the server side so that internal
// details are not exposed to the caller.  Additionally, a request payload that is not well-formed will result in a 400
// http status.
func NewMiddleware(ctx *gin.Context) {
	ctx.Next()
	for _, err := range ctx.Errors {
		if err.Error() == BadRequest {
			handleBadRequest(ctx)
			return
		}
		var businessError *business.Error
		if errors.As(err, &businessError) {
			handleBusinessError(ctx, businessError)
			return
		}
	}
	if len(ctx.Errors) > 0 {
		handleSystemError(ctx)
	}
}

func handleBadRequest(ctx *gin.Context) {
	response := &ErrorResponse{
		Message: badRequestErrorMessage,
	}
	ctx.JSON(http.StatusBadRequest, response)
}

func handleBusinessError(ctx *gin.Context, businessError *business.Error) {
	ctx.JSON(http.StatusUnprocessableEntity, businessError)
}

func handleSystemError(ctx *gin.Context) {
	log.Printf("system error: %v\n", ctx.Errors)
	response := &ErrorResponse{
		Message: systemErrorMessage,
	}
	ctx.JSON(http.StatusInternalServerError, response)
}
