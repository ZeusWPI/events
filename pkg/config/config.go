// Package config provides all configuration related functions
package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func bindEnv(key string) {
	envName := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	// nolint:errcheck // we do not care if it can get binded
	viper.BindEnv(key, envName)
}

// Init initializes the configuration
func Init() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	viper.AutomaticEnv()
	env := GetDefaultString("app.env", "development")
	env = strings.ToLower(env)

	viper.SetConfigName(env + ".yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	return viper.ReadInConfig()
}

// GetString returns the value of the key in string
func GetString(key string) string {
	bindEnv(key)
	return viper.GetString(key)
}

// GetDefaultString returns the value of the key in string or a default value
func GetDefaultString(key, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return GetString(key)
}

// GetInt returns the value of the key in int
func GetInt(key string) int {
	bindEnv(key)
	return viper.GetInt(key)
}

// GetDefaultInt returns the value of the key in int or a default value
func GetDefaultInt(key string, defaultVal int) int {
	viper.SetDefault(key, defaultVal)
	return GetInt(key)
}

// GetUint16 returns the value of the key in uint16
func GetUint16(key string) uint16 {
	bindEnv(key)
	return viper.GetUint16(key)
}

// GetDefaultUint16 returns the value of the key in uint16 or a default value
func GetDefaultUint16(key string, defaultVal uint16) uint16 {
	viper.SetDefault(key, defaultVal)
	return GetUint16(key)
}

// GetBool returns the value of the key in bool
func GetBool(key string) bool {
	bindEnv(key)
	return viper.GetBool(key)
}

// GetDefaultBool returns the value of the key in bool or a default value
func GetDefaultBool(key string, defaultVal bool) bool {
	viper.SetDefault(key, defaultVal)
	return GetBool(key)
}
