package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os"
	"strings"

	fire "firebase.google.com/go"
	"github.com/labstack/echo"
	"google.golang.org/api/option"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8081"},
		AllowHeaders: []string{echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET},
	}))

	// router
	e.GET("/public", public)
	e.GET("/private", private)

	// server start
	e.Logger.Fatal(e.Start(":8000"))
}

func public(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	return c.JSON(http.StatusOK, "public is success.")
}
func private(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	opt := option.WithCredentialsFile("/Users/yoshinori.hisakawa/keys/firebase-key.json")
	app, err := fire.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	auth, err := app.Auth(context.Background())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	// クライアントから送られてきた JWT 取得
	authHeader := c.Request().Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)

	// JWT の検証
	token, err := auth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		u := fmt.Sprintf("error verifying ID token: %v\n", err)
		return c.JSON(http.StatusBadRequest, u)
	}
	uid := token.Claims["user_id"]

	u := fmt.Sprintf("public is success. uid is %s", uid)

	return c.JSON(http.StatusOK, u)
}
