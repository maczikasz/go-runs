// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package auth

import (
	"github.com/maczikasz/go-runs/internal/auth"
	"github.com/maczikasz/go-runs/internal/wire/auth/config"
)

// Injectors from wire.go:

func InitializeAuthContext() *auth.AuthContext {
	oauth2Config := config.NewAuthConfig()
	logoutUrl := config.NewLogoutUrl()
	authContext := auth.NewAuthContext(oauth2Config, logoutUrl)
	return authContext
}
