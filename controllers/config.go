package controllers

import (
	"github.com/spf13/viper"
	"log"
)

func ViperConfigVariable(key string) string {
	// type of config file
	viper.SetConfigType("yaml")
	// name of config file (without extension)
	viper.SetConfigName("config")
	// look for config in the working directory
	viper.AddConfigPath(".")

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// viper.Get() returns an empty interface{}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}
