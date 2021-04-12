package main

import (
	"github.com/maczikasz/go-runs/internal/auth"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"sync"
)

func main() {
	log.SetLevel(log.TraceLevel)
	wg := sync.WaitGroup{}
	wg.Add(1)

	fusionAuthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:3000/oauth-callback",
		ClientID:     "3b37d069-f806-465e-92a2-3ef9f8ef9ca3",
		ClientSecret: "47_bqYSB0azC8hSWyUrCTpL6StVRPfAjBVp5Kv5Tz6k",
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "http://localhost:9011/oauth2/authorize",
			TokenURL:  "http://localhost:9011/oauth2/token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}

	authContext := auth.NewAuthContext(fusionAuthConfig," http://localhost:9011/oauth2/logout")

	auth.StartHttpServer(&wg, authContext)
	wg.Wait()

}
