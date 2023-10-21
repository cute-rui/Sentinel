package config

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var Conf *viper.Viper

func init() {
	Conf = viper.New()

	Conf.SetConfigType("toml")
	Conf.SetConfigName("sentinel")
	Conf.AddConfigPath(`.`)

	Conf.SetDefault("HTTP.ListenAddr", "")

	Conf.SetDefault("Database.User", "")
	Conf.SetDefault("Database.Pass", "")
	Conf.SetDefault("Database.Host", "")
	Conf.SetDefault("Database.Port", 0)
	Conf.SetDefault("Database.Name", "")

	Conf.SetDefault("JWT.AccessSecret", "")
	Conf.SetDefault("Admin.Secret", "")
	Conf.SetDefault("Verify.Interval", 300)

	Conf.SetDefault("OIDC.32BKey", "")
	Conf.SetDefault("OIDC.WebID", "")
	Conf.SetDefault("OIDC.WebSecret", "")
	Conf.SetDefault("OIDC.RedirectAllowedList", []string{})
	Conf.SetDefault("OIDC.Issuer", "")

	replacer := strings.NewReplacer(".", "_")
	Conf.SetEnvKeyReplacer(replacer)
	err := Conf.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		_, err := os.Create("./sentinel.toml")
		if err != nil {
			panic(err)
		}

		err = Conf.WriteConfig()
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	Conf.WatchConfig()
	Conf.OnConfigChange(func(in fsnotify.Event) {
		err := Conf.ReadInConfig()
		if err != nil {
			log.Println(err)
		}
	})
}
