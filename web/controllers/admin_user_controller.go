package controllers

import (
	"go_url_shortener/web/database"
	"go_url_shortener/web/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func AdminShowUsers(c echo.Context) error {
	var users []models.UserPresentation

	database.DB.Table("users").Where("is_admin = ?", false).Find(&users)

	return c.Render(http.StatusOK, "admin_show_users.html", echo.Map{
		"IsAdmin": true,
		"IsUser":  false,
		"users":   users,
	})

}

func AdminUserUpdateGET(c echo.Context) error {
	user_id := c.Param("user_id")

	var user models.UserPresentation

	database.DB.Table("users").Where("id = ?", user_id).First(&user)
	return c.Render(http.StatusOK, "admin_user_update.html", echo.Map{
		"Name":    user.Name,
		"Email":   user.Email,
		"IsAdmin": user.IsAdmin,
		"IsUser":  !user.IsAdmin,
	})
}

func AdminUserUpdatePUT(c echo.Context) error {
	user_id := c.Param("user_id")
	var user models.UserRegisterForm

	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	id, err := strconv.ParseUint(user_id, 10, 32)

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
	}
	database.DB.Updates(&updatedUser)
	return c.Redirect(http.StatusSeeOther, "/yonetici/kullanicilar")
}

func AdminUserDelete(c echo.Context) error {
	user_id := c.Param("user_id")

	database.DB.Delete(&models.User{}, user_id)
	return c.Redirect(http.StatusSeeOther, "/yonetici/kullanicilar")
}
