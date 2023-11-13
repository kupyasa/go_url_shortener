package routes

import (
	"go_url_shortener/web/controllers"
	custom_middlewares "go_url_shortener/web/middlewares"

	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo) {
	e.GET("/", controllers.Index, custom_middlewares.User)
	e.GET("/giris", controllers.Login, custom_middlewares.Guest)
	e.GET("/kayit", controllers.Register, custom_middlewares.Guest)
	e.POST("/kayit", controllers.RegisterPOST)
	e.POST("/giris", controllers.LoginPOST)
	e.GET("/kullanici", controllers.UserGET, custom_middlewares.User)
	e.GET("/kullanici/guncelle", controllers.UserUpdateGET, custom_middlewares.User)
	e.POST("/kullanici/guncelle", controllers.UserUpdatePUT, custom_middlewares.User)
	e.GET("/yonetici/kullanicilar", controllers.AdminShowUsers, custom_middlewares.Admin)
	e.GET("/yonetici/kullanici/:user_id/duzenle", controllers.AdminUserUpdateGET, custom_middlewares.Admin)
	e.POST("/yonetici/kullanici/:user_id/duzenle", controllers.AdminUserUpdatePUT, custom_middlewares.Admin)
	e.POST("/yonetici/kullanici/:user_id/sil", controllers.AdminUserDelete, custom_middlewares.Admin)
	e.POST("/linkolustur", controllers.CreateShortLink, custom_middlewares.User)
	e.GET("/link/:hashed_string", controllers.LinkClicked, custom_middlewares.User)
	e.GET("/linklerim", controllers.GetMyLinks, custom_middlewares.User)
	e.POST("/linklerim/:link_id/sil", controllers.UserLinkDelete, custom_middlewares.User)
	e.GET("/yonetici/linkler", controllers.AdminGetLinks, custom_middlewares.User)
	e.POST("/yonetici/linkler/:link_id/sil", controllers.AdminDeleteLink, custom_middlewares.User)
	e.POST("/cikis", controllers.Logout)

}
