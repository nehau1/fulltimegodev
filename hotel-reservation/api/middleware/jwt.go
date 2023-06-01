package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--- JWT middleware")
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	claims, err := validateToken(token); 
	if err != nil {
		return err
	}
	fmt.Println("claims: ", claims)
	expires := claims["expires"].(float64)
	if time.Now().Unix() > int64(expires) {
		return fmt.Errorf("token expired")
	}
	fmt.Println("expires: ", expires)
	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims,error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		// secret := os.Getenv("JWT_SECRET")
		return []byte("secret"), nil
	})

	if err != nil {
		fmt.Println("failed to parse token: ", err)
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok{
		return claims, nil
	}

	return nil,fmt.Errorf("unauthorized")
}
