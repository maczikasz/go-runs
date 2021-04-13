package auth

import "golang.org/x/oauth2"

type AuthContext struct {
	oauthConfig *oauth2.Config
	logoutUrl   string
}

func (c AuthContext) LogoutUrl() string {
	return c.logoutUrl + "?client_id=" + c.oauthConfig.ClientID
}

type LogoutUrl string

func NewAuthContext(oauthConfig *oauth2.Config, logoutUrl LogoutUrl) *AuthContext {
	return &AuthContext{oauthConfig: oauthConfig, logoutUrl: string(logoutUrl)}
}
