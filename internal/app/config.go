package app

import (
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DevSearch      string   `validate:"required"`
	Token          string   `validate:"required"`
	Timeout        int      `validate:"required"`
	AllowedIDs     []string `validate:"required"`
	FTPuser        string   `validate:"required"`
	PathToDownload string   `validate:"required"`
	PathToDB       string   `validate:"required"`
}

func loadConfig(p string) (*Config, error) {
	fullName := strings.Split(strings.Split(p, "/")[len(strings.Split(p, "/"))-1], ".")

	viper.SetConfigType(fullName[1]) // Type
	viper.SetConfigName(fullName[0]) // Filename without type

	pathStr := strings.Builder{}
	for _, v := range strings.Split(p, "/")[:len(strings.Split(p, "/"))-1] {
		pathStr.WriteString(v)
		pathStr.WriteString("/")
	}

	viper.AddConfigPath(pathStr.String())

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("LoadConfig: %w", err)
	}

	conf := &Config{}

	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, fmt.Errorf("LoadConfig(Unmarshal): %w", err)
	}

	if err := valid(conf); err != nil {
		return nil, fmt.Errorf("LoadConfig(Validate): %w", err)
	}

	return conf, nil
}

func valid(conf *Config) error {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(conf); err != nil {
		return fmt.Errorf("validation faild: %w", err)
	}

	return nil
}
