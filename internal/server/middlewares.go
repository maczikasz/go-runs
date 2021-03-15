package server

import (
	"net/http"
	"net/url"
)

type CORSPreflightOriginMiddleware struct {
	AcceptedOrigins map[string]bool
}

const (
	allowOriginHeaderName = "Access-Control-Allow-Origin"
	allowHeaderHeaderName = "Access-Control-Allow-Headers"
)

func (receiver CORSPreflightOriginMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if origin := r.Header.Get("origin"); origin != "" {
			parse, _ := url.Parse(origin)
			if parse != nil {
				if receiver.AcceptedOrigins[parse.Host] == true {
					w.Header().Add(allowOriginHeaderName, origin)
					//TODO fix this
					w.Header().Add(allowHeaderHeaderName, "*")
				}
			}
		}
		next.ServeHTTP(w, r)
	},
	)
}
