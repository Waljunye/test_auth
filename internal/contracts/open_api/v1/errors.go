package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HttpError struct {
	Code      int
	Message   string
	BaseError error
	Body      interface{}
}

func (he *HttpError) Error() string {
	if he.BaseError != nil {
		return errors.Wrap(he.BaseError, he.Message).Error()
	}
	return he.Message
}

func NewHttpError(code int, message string, data interface{}, err error) *HttpError {
	return &HttpError{
		Code:      code,
		Message:   message,
		Body:      data,
		BaseError: err,
	}
}
func processError(ctx *gin.Context, err *HttpError) {
	if err == nil {
		return
	}
	ctx.JSON(err.Code, struct {
		Code    int
		Message string
	}{
		Code:    err.Code,
		Message: err.Error(),
	})

	return
}
