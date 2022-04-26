package httperr

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpErr struct {
	statusCode int
	message    string
}

func (e HttpErr) Error() string {
	return e.message
}

func Handle(c *gin.Context, err error) {
	if hErr, ok := err.(HttpErr); ok { // If it is httperr - get status code from it
		c.JSON(hErr.statusCode, gin.H{
			"error": hErr.message,
		})
	} else { // Otherwise use internal error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
}

func New(msg string, code int) error {
	return HttpErr{
		message:    msg,
		statusCode: code,
	}
}

func Wrap(err error, msg string) error {
	return errors.New(fmt.Sprintf("%s: %v", msg, err))
}

func WrapHttp(err error, msg string, code int) error {
	return HttpErr{
		message:    fmt.Sprintf("%s: %v", msg, err),
		statusCode: code,
	}
}

func ValidationError(field, msg string) error {
	return HttpErr{
		message:    fmt.Sprintf("%s: %s", field, msg),
		statusCode: http.StatusBadRequest,
	}
}
