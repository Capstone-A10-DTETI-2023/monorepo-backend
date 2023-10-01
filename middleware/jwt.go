package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	_ "github.com/joho/godotenv/autoload"
)

type Claims struct {
	ID      uint
	Role_ID uint
	jwt.StandardClaims
}

func GenerateJWT (id uint, role_id uint) (string, error) {
	claims := Claims{
		ID:      id,
		Role_ID: role_id,
		StandardClaims: jwt.StandardClaims{},
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Hour * 72).Unix()
	claims.Issuer = os.Getenv("JWT_ISSUER")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ParseJWT (token string) (*Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(token, &claims, func (t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func IsAuthenticated (c *fiber.Ctx) error {
	token := c.Cookies("token")
	claims, err := ParseJWT(token)

	if err != nil || claims.ExpiresAt <= 0 || claims.Issuer != os.Getenv("JWT_ISSUER") {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}

func IsAdmin (c *fiber.Ctx) bool {
	token := c.Cookies("token")
	claims, err := ParseJWT(token)

	if err != nil || claims.ExpiresAt <= 0 || claims.Role_ID != 1 {
		return false
	}

	return true
}
