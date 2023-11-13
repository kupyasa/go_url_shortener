package main

import (
	"errors"
	"go_url_shortener/web/database"
	"go_url_shortener/web/routes"
	"go_url_shortener/web/utils"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}

func main() {
	utils.SetupViper()
	database.Connect()
	templates := make(map[string]*template.Template)
	templates["index.html"] = template.Must(template.ParseFiles("./ui/html/index.html", "./ui/html/navbar.html", "./ui/html/base.html"))
	templates["login.html"] = template.Must(template.ParseFiles("./ui/html/login.html", "./ui/html/navbar.html", "./ui/html/base.html"))
	templates["register.html"] = template.Must(template.ParseFiles("./ui/html/register.html", "./ui/html/navbar.html", "./ui/html/base.html"))
	templates["profile.html"] = template.Must(template.ParseFiles("./ui/html/profile.html", "./ui/html/navbar.html", "./ui/html/base.html"))
	templates["profile_update.html"] = template.Must(template.ParseFiles("./ui/html/profile_update.html", "./ui/html/navbar.html", "./ui/html/base.html"))
	templates["admin_show_users.html"] = template.Must(template.ParseFiles("./ui/html/admin_show_users.html", "./ui/html/navbar.html", "./ui/html/base.html"))
	templates["admin_user_update.html"] = template.Must(template.ParseFiles("./ui/html/admin_user_update.html", "./ui/html/navbar.html", "./ui/html/base.html"))
	templates["my_links.html"] = template.Must(template.ParseFiles("./ui/html/my_links.html", "./ui/html/navbar.html", "./ui/html/base.html"))
	templates["admin_links.html"] = template.Must(template.ParseFiles("./ui/html/admin_links.html", "./ui/html/navbar.html", "./ui/html/base.html"))
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
	}))
	e.Static("/static", "./ui/static")
	e.Renderer = &TemplateRegistry{
		templates: templates,
	}
	routes.Setup(e)
	e.Logger.Fatal(e.Start(":1323"))
}
