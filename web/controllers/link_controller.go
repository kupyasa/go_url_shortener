package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go_url_shortener/web/database"
	"go_url_shortener/web/models"
	"go_url_shortener/web/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
)

func CreateShortLink(c echo.Context) error {
	orginal_url := c.FormValue("original_url")
	tag := strings.TrimSpace(c.FormValue("tag"))
	expiration_date := c.FormValue("expiration_date")
	if !govalidator.IsURL(orginal_url) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid URL",
		})
	}
	//return c.String(http.StatusAccepted, expiration_date)
	expires_at, _ := time.ParseInLocation("2006-01-02T15:04", expiration_date, time.Local)
	if len(expiration_date) == 0 {
		expires_at = time.Now().AddDate(1, 0, 0)
	}

	cookie, err := c.Cookie("jwt")

	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"message": "JWT Not Found"})
	}

	claims, err := utils.GetClaims(cookie)

	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"message": "Could not parse claims"})
	}

	if len(tag) == 0 {
		tag = gonanoid.Must(12)
	}

	var tag_check models.Link

	for {
		database.DB.Where("tag = ?", tag).First(&tag_check)

		if tag_check.Id != 0 {
			tag = gonanoid.Must(12)
		} else {
			break
		}
	}

	hasher := sha256.New()

	hasher.Write([]byte(tag))

	hashSum := hasher.Sum(nil)

	shortened_url := hex.EncodeToString(hashSum)

	id, _ := strconv.ParseUint(claims.Issuer, 10, 32)

	short_link := models.Link{
		OriginalUrl:  orginal_url,
		ShortenedUrl: viper.GetString("BASE_URL") + "/link/" + shortened_url,
		UserId:       uint(id),
		Tag:          tag,
		CreatedAt:    time.Now().UTC(),
		ExpiresAt:    expires_at.UTC(),
		ClickCount:   0,
	}

	database.DB.Create(&short_link)

	html := fmt.Sprintf(`
	<div
	class="card bg-base-100 shadow-xl m-8"
	id="short-url-card"
  >
	<div class="card-body items-center text-center">
	  <h2 class="card-title">Yeni URL Oluşturuldu</h2>
	  <p>
		Orijinal URL :</p>
		<a
		  target="_blank"
		  href="%s"
		  >%s</a
		>
	  <p>
		Kısaltılmış URL :</p>
		<a
		  target="_blank"
		  href="%s"
		  >%s</a
		>
	  <p>Özel Ad : %s</p>
	  <p>
		Oluşturulduğu Tarih :
		%s
	  </p>
	  <p>
		Son Geçerlilik Tarih :
		%s
	  </p>
	  <div class="card-actions justify-end my-4">
		<button class="btn btn-primary" hx-get="/" hx-target="#short-url-card"
		hx-swap="delete">
		  Kapat
		</button>
	  </div>
	</div>
  </div>
	`, orginal_url, orginal_url, short_link.ShortenedUrl, short_link.ShortenedUrl, short_link.Tag, time.Now().Format("02-01-2006 15:04"), expires_at.Format("02-01-2006 15:04"))

	return c.HTML(http.StatusOK, html)
}

func LinkClicked(c echo.Context) error {
	hashed_string := c.Param("hashed_string")

	var link models.Link

	database.DB.Where("shortened_url LIKE ?", "%"+hashed_string+"%").First(&link)

	if link.Id == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Link Not Found",
		})
	}

	if time.Now().Compare(link.ExpiresAt) == -1 {
		database.DB.Model(&link).Update("click_count", link.ClickCount+1)
		return c.Redirect(http.StatusSeeOther, link.OriginalUrl)
	} else if time.Now().Compare(link.ExpiresAt) == 1 {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Link is expired",
		})
	}

	return c.JSON(http.StatusInternalServerError, echo.Map{
		"message": "Couldn't redirect to URL",
	})
}

func GetMyLinks(c echo.Context) error {
	cookie, err := c.Cookie("jwt")

	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"message": "JWT Not Found"})
	}

	claims, err := utils.GetClaims(cookie)

	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"message": "Could not parse claims"})

	}

	var links []models.Link

	database.DB.Where("user_id = ?", claims.Issuer).Find(&links)

	return c.Render(http.StatusOK, "my_links.html", echo.Map{
		"IsAdmin": claims.Admin,
		"IsUser":  !claims.Admin,
		"links":   links})
}

func UserLinkDelete(c echo.Context) error {
	link_id := c.Param("link_id")

	database.DB.Delete(&models.Link{}, link_id)
	return c.Redirect(http.StatusSeeOther, "/linklerim")
}

func AdminGetLinks(c echo.Context) error {

	var links []models.Link

	database.DB.Find(&links)

	return c.Render(http.StatusOK, "admin_links.html", echo.Map{
		"IsAdmin": true,
		"IsUser":  false,
		"links":   links})
}

func AdminDeleteLink(c echo.Context) error {
	link_id := c.Param("link_id")

	database.DB.Delete(&models.Link{}, link_id)
	return c.Redirect(http.StatusSeeOther, "/yonetici/linkler")
}
