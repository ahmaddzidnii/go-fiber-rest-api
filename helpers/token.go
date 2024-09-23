package helpers

import (
	"time"

	"github.com/ahmaddzidnii/go-fiber-rest-api/models"
	"github.com/golang-jwt/jwt/v4"
)

var SecretKey = []byte("secret");

type MyJwtClaims struct {
	Id uint `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User,expiredDate *jwt.NumericDate) (string, error) {
	claims := MyJwtClaims{
		user.Id,
		user.FullName,
		user.Username,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: expiredDate,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims);

	ss,err := token.SignedString(SecretKey);

	return ss,err;
}

func ClaimJWT(token string) (*MyJwtClaims, error) {
	tkn, err := jwt.ParseWithClaims(token, &MyJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil;
	})

	if err != nil {
		return nil, err;
	}

	claims, ok := tkn.Claims.(*MyJwtClaims);

	if !ok || !tkn.Valid {
		return nil, err;
	}

	return claims, nil;
}

