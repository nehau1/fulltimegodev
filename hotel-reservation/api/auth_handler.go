package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/Stiffjobs/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{userStore: userStore}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResp struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type:    "error",
		Message: "invalid credentials",
	})
}

// A handler should only do:
//  - serialization of the incoming request(JSON)
//	- do some data fetching from the db
//  - call some business logic
//  - return a response

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err

	}
	user, err := h.userStore.GetByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}
	if valid := types.IsValidPassword(user.EncryptedPassword, params.Password); !valid {
		return invalidCredentials(c)
	}
	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 24 * 7).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println("failed to sign token: ", err)
	}
	return tokenStr
}
