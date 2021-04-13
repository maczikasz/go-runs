package auth

import (
	gocontext "context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func StartHttpServer(wg *sync.WaitGroup, authContext *AuthContext) {
	defer wg.Done()
	r := SetupRouter(authContext, []string{"http://localhost:3000"})

	//TODO error handle
	_ = r.Run()
}

func SetupRouter(authContext *AuthContext, acceptedOrigins []string) *gin.Engine {
	r := gin.Default()

	config := cors.Config{
		AllowHeaders: []string{"Content-type", "Origin", "Authorization"},
		AllowMethods: []string{"POST", "GET"},
	}
	if len(acceptedOrigins) == 0 {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = acceptedOrigins
	}

	r.Use(cors.New(config))

	r.GET("/login", redirectOnLoginHandler(authContext))
	r.GET("/logout", redirectOnLogoutHandler(authContext))
	r.GET("/user", fetchUserDataHandler())
	r.POST("/token", fetchTokenHadler(authContext))

	return r
}

func fetchTokenHadler(authContext *AuthContext) func(context *gin.Context) {
	return func(context *gin.Context) {
		type Code struct {
			Code string `json:"code"`
		}
		var code Code
		_ = context.BindJSON(&code)

		token, _ := authContext.oauthConfig.Exchange(gocontext.Background(), code.Code)

		fmt.Printf("Token %#v", token)

		context.JSON(200, token)
	}
}

func fetchUserDataHandler() func(context *gin.Context) {
	return func(context *gin.Context) {
		request, _ := http.NewRequest(http.MethodGet, "http://localhost:9011/oauth2/userinfo", nil)

		request.Header.Add("Authorization", context.GetHeader("Authorization"))

		response, _ := http.DefaultClient.Do(request)
		context.DataFromReader(200, response.ContentLength, response.Header.Get("Content-type"), response.Body, nil)

	}
}

func redirectOnLoginHandler(authContext *AuthContext) func(context *gin.Context) {
	return func(context *gin.Context) {
		//TODO change this
		redirectUrl := authContext.oauthConfig.AuthCodeURL("test")
		context.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	}
}

func redirectOnLogoutHandler(authContext *AuthContext) func(context *gin.Context) {
	return func(context *gin.Context) {
		redirectUrl := authContext.LogoutUrl()
		context.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	}
}
