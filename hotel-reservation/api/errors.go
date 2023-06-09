package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiErr, ok := err.(Error); ok {
		return c.Status(apiErr.Code).JSON(apiErr)
	}
	apiErr := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiErr.Code).JSON(apiErr)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Err:  msg,
	}
}

func ErrUnauthorized() Error {
	return NewError(http.StatusUnauthorized, "unauthorized request")
}
func ErrBadRequest() Error {
	return NewError(http.StatusBadRequest, "bad request")
}

func ErrNotFound() Error {
	return NewError(http.StatusNotFound, "not found")
}

func ErrResourceNotFound(res string) Error {
	return NewError(http.StatusNotFound, fmt.Sprintf("%s not found", res))
}

func ErrInvalidID() Error {
	return NewError(http.StatusNotFound, "invalid id given")
}
