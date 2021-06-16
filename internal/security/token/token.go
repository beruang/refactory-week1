package token

import (
	"github.com/dgrijalva/jwt-go"
	"refactory/notes/internal/app/model"
	"time"
)

type Token struct {
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	RoleId   int    `json:"role_id"`
}

func GenerateToken(session model.Session) (string, error) {
	claims := Token{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
		UserId:   session.UserId,
		Username: session.Username,
		RoleId:   session.RoleId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}
