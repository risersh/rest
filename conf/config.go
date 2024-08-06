package conf

import (
	"log"
	"os"

	"github.com/risersh/config/config"
	"github.com/risersh/notifications"
)

var Config *Conf

type Conf struct {
	config.BaseConfig
	Server struct {
		Port int `yaml:"port" env:"PORT" env-default:"8080"`
		Cors struct {
			Origins []string `yaml:"origins" env:"CORS_ORIGINS"`
		} `yaml:"cors" env-prefix:"CORS_"`
	} `yaml:"server" env-prefix:"SERVER_"`
	Session struct {
		Duration int `yaml:"duration" env:"SESSIONS_DURATION" env-default:"24"`
	} `yaml:"session" env-prefix:"SESSIONS_"`
	Branding     notifications.Branding `yaml:"branding"`
	Registration struct {
		Code string `yaml:"code" env:"REGISTRATION_INTERNAL_CODE"`
	} `yaml:"registration" env-prefix:"REGISTRATION_"`
}

func Init() {
	var err error
	Config, err = config.GetConfig[Conf](config.Environment(os.Getenv("ENV")))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	os.Setenv("DATABASE_URI", Config.Database.URI)
}
