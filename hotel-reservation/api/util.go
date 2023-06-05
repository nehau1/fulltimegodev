package api

import (
	"github.com/Stiffjobs/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthedUser(c *fiber.Ctx) (*types.User, error) {

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}
	return user, nil
}
