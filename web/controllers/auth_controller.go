package controllers

import (
	"go_url_shortener/web/database"
	"go_url_shortener/web/models"
	"go_url_shortener/web/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func Index(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "JWT value is missing . Please login",
		})
	}

	claims, err := utils.GetClaims(cookie)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Unauthorized Access",
		})
	}

	return c.Render(http.StatusOK, "index.html", echo.Map{
		"IsAdmin": claims.Admin,
		"IsUser":  !claims.Admin,
	})
}

func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", echo.Map{
		"IsAdmin": false,
		"IsUser":  false,
	})
}

func Register(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", echo.Map{
		"IsAdmin": false,
		"IsUser":  false,
	})
}

func RegisterPOST(c echo.Context) error {
	var user models.UserRegisterForm

	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(password),
		IsAdmin:  false,
	}
	database.DB.Create(&newUser)
	return c.Redirect(http.StatusSeeOther, "/giris")
}

func LoginPOST(c echo.Context) error {
	var data = make(map[string]string)

	err := c.Bind(&data)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Incorrect password",
		})
	}

	//return c.JSON(http.StatusOK, user)

	claims := &models.JWTCustomClaims{
		Name:  user.Name,
		Admin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			Issuer:    strconv.Itoa(int(user.Id)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Token signing error",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.Expires = time.Now().Add(time.Hour * 72)
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/")
}

func Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-(time.Minute))
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/giris")
}
