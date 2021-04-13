package config

import (
	"github.com/maczikasz/go-runs/internal/auth"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func NewAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  viper.GetString("auth.redirect"),
		ClientID:     viper.GetString("auth.client.id"),
		ClientSecret: viper.GetString("auth.client.secret"),
		Scopes:       viper.GetStringSlice("auth.scopes"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   viper.GetString("auth.endpoint.auth"),
			TokenURL:  viper.GetString("auth.endpoint.token"),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
}

func NewLogoutUrl() auth.LogoutUrl {
	return auth.LogoutUrl(viper.GetString("auth.endpoint.logout"))
}
