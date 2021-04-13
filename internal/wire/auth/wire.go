//+build wireinject

package auth

import (
	"github.com/google/wire"
	"github.com/maczikasz/go-runs/internal/auth"
	"github.com/maczikasz/go-runs/internal/wire/auth/config"
)

func InitializeAuthContext() *auth.AuthContext {
	wire.Build(auth.NewAuthContext, config.NewAuthConfig, config.NewLogoutUrl)

	return &auth.AuthContext{}
}
