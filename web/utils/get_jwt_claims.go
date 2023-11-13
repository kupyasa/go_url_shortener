package utils

import (
	"go_url_shortener/web/models"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func GetClaims(cookie *http.Cookie) (*models.JWTCustomClaims, error) {

	token, err := jwt.ParseWithClaims(cookie.Value, &models.JWTCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT_SECRET")), nil
	})

	if err != nil {
		println(err.Error())
		return nil, err
	}

	return token.Claims.(*models.JWTCustomClaims), nil
}
