package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rrawrriw/angular-sauth-handler"
	"github.com/rrawrriw/sj-handler"
)

const (
	SessionColl = "Session"
)

func main() {
	app, err := sj.NewApp("SJ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(app.Specs)

	singIn := aauth.AngularSignIn(findUser(app))
	newUser := sj.NewAppHandler(sj.NewUserHandler, app)
	//singedIn := aauth.AngularAuth(app.DB(), SessionColl)

	host := app.Specs.Host
	port := app.Specs.Port
	srvRes := host + ":" + strconv.Itoa(port)

	publicDir := app.Specs.PublicDir
	htmlDir := path.Join(publicDir, "html")

	srv := gin.Default()
	srv.Use(sj.Serve("/", sj.LocalFile(htmlDir, false)))
	srv.Static("/public", publicDir)
	srv.GET("/SignIn/:name/:pass", singIn)
	srv.POST("/User", newUser)

	srv.Run(srvRes)
}

func findUser(app sj.AppContext) aauth.FindUser {
	return func(name string) (aauth.User, error) {
		return sj.FindUser(app.DB(), name)
	}
}
