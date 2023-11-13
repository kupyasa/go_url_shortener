package database

import (
	"go_url_shortener/web/models"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	println(viper.GetString("DATABASE_CONN"))
	connection, err := gorm.Open(postgres.Open(viper.GetString("DATABASE_CONN")), &gorm.Config{})

	if err != nil {
		panic("Could not connect database")
	}
	DB = connection
	connection.AutoMigrate(&models.User{}, &models.Link{})
	println("Connected")
}
