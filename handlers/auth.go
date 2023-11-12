package handlers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"strconv"
)

type Identity struct {
	Id    int
	Email string
}

var IdentityStruct *Identity

var Secret string = "wow secret is here"
var Aud string = "wow aud is here"

func WithAuth() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			t := ctx.Request().Header.Get("Authorization")
			if t == "" {
				return ctx.JSON(401, "you are not authorized by server")
			}
			t = t[len("Bearer "):]
			token, err := validateToken(t)
			if err != nil {
				return ctx.JSON(401, "you are not authorized by server")
			}
			ident, err := newIdentity(token)
			if err != nil {
				return ctx.JSON(401, "you are not authorized by server")
			}

			IdentityStruct = ident
			return h(ctx)
		}
	}
}

func newIdentity(t *jwt.Token) (*Identity, error) {
	claims := t.Claims.(jwt.MapClaims)
	id, err := strconv.Atoi(claims["jti"].(string))
	if err != nil {
		return nil, err
	}

	return &Identity{
		Id:    id,
		Email: claims["sub"].(string),
	}, nil
}

func NewJWTToken(claims jwt.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(Secret))
}

func validateToken(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid jwt token")
	}
	return t, nil
}
