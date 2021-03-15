package infra

import (
	"github.com/gorilla/mux"
	"github.com/maczikasz/go-runs/internal/util"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Return200Handler struct {
}

func (r Return200Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
}

func TestCorrectCORSHeaders(t *testing.T) {

	Convey("Given http router with CORS middleware", t, func() {
		router := mux.NewRouter()
		router.Handle("/test", Return200Handler{})
		middleware := CORSPreflightOriginMiddleware{AcceptedOrigins: util.ToSet([]string{"test:1234"})}
		router.Use(middleware.Middleware)

		Convey("When Origin is specified and is http", func() {
			r, _ := http.NewRequest(http.MethodPost, "/test", nil)
			w := httptest.NewRecorder()

			origin := "http://test:1234"
			r.Header.Add("Origin", origin)

			router.ServeHTTP(w, r)
			Convey("Then CORS preflight headers are present", func() {
				So(w.Header().Get("Access-Control-Allow-Origin"), ShouldEqual, origin)
				So(w.Header().Get("Access-Control-Allow-Headers"), ShouldEqual, "*")
			})
		})

		Convey("When Origin is specified and is https", func() {
			r, _ := http.NewRequest(http.MethodPost, "/test", nil)
			w := httptest.NewRecorder()

			origin := "https://test:1234"
			r.Header.Add("Origin", origin)

			router.ServeHTTP(w, r)
			Convey("Then CORS preflight headers are present", func() {
				So(w.Header().Get("Access-Control-Allow-Origin"), ShouldEqual, origin)
				So(w.Header().Get("Access-Control-Allow-Headers"), ShouldEqual, "*")
			})
		})
	})
}
