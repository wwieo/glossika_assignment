package errorx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type errCode int

type ErrorResponse struct {
	Code    errCode `json:"code"`
	Message string  `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func newErrorResponse(code errCode, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

func RespondWithError(ctx *gin.Context, statusCode int, err error) {
	if statusCode == http.StatusInternalServerError {
		ctx.AbortWithStatusJSON(statusCode, err.Error())
		return
	}

	var errorResponse ErrorResponse
	if respErr, ok := err.(ErrorResponse); ok {
		errorResponse = respErr
	} else {
		errMsg := fmt.Sprintf("service internal error: %s", err.Error())
		errorResponse = newErrorResponse(serviceInternal, errMsg)
	}

	ctx.AbortWithStatusJSON(statusCode, errorResponse)
}
