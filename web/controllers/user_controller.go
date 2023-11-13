package controllers

import (
	"go_url_shortener/web/database"
	"go_url_shortener/web/models"
	"go_url_shortener/web/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func UserGET(c echo.Context) error {
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

	var user models.UserPresentation

	database.DB.Table("users").Where("id = ?", claims.Issuer).First(&user)
	return c.Render(http.StatusOK, "profile.html", echo.Map{
		"Name":    user.Name,
		"Email":   user.Email,
		"IsAdmin": user.IsAdmin,
		"IsUser":  !user.IsAdmin,
	})
}

func UserUpdateGET(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Unauthenticated Access . Please login",
		})
	}

	claims, err := utils.GetClaims(cookie)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Unauthorized Access",
		})
	}

	var user models.UserPresentation

	database.DB.Table("users").Where("id = ?", claims.Issuer).First(&user)
	return c.Render(http.StatusOK, "profile_update.html", echo.Map{
		"Name":    user.Name,
		"Email":   user.Email,
		"IsAdmin": user.IsAdmin,
		"IsUser":  !user.IsAdmin,
	})
}

func UserUpdatePUT(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Unauthenticated Access . Please login",
		})
	}
	claims, err := utils.GetClaims(cookie)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Unauthorized Access",
		})
	}

	var user models.UserRegisterForm

	err = c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	id, err := strconv.ParseUint(claims.Issuer, 10, 32)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Couldn't convert id",
		})
	}

	updatedUser := models.User{
		Id:       uint(id),
		Name:     user.Name,
		Email:    user.Email,
		Password: string(password),
		IsAdmin:  claims.Admin,
	}
	database.DB.Save(&updatedUser)
	return c.Redirect(http.StatusSeeOther, "/kullanici")
}
