package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func SetupViper() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Read in the environment variables from the .env file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading .env file: %s\n", err)
		os.Exit(1)
	}
}
