package custom_middlewares

import (
	"go_url_shortener/web/database"
	"go_url_shortener/web/models"
	"go_url_shortener/web/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

var User = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			println(err.Error())
			return c.Redirect(http.StatusSeeOther, "/giris")
		}
		claims, err := utils.GetClaims(cookie)

		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/giris")
		}

		var user models.User

		database.DB.Where("id = ?", claims.Issuer).First(&user)

		if user.Id == 0 {
			return c.Redirect(http.StatusSeeOther, "/giris")
		}

		return next(c)
	}
}

var Admin = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}
		claims, err := utils.GetClaims(cookie)

		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		var user models.User

		database.DB.Where("id = ?", claims.Issuer).First(&user)

		if !user.IsAdmin {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		return next(c)
	}
}

var Guest = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			println(err.Error())
			return next(c)
		}
		claims, err := utils.GetClaims(cookie)

		if err != nil {
			println(err.Error())
			return next(c)
		}

		var user models.User

		database.DB.Raw("SELECT id,name,email FROM users WHERE id = ?", claims.Issuer).First(&user)

		if user.Id != 0 {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		return next(c)
	}
}
