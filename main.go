package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

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

	sessionColl := app.DB().C(SessionColl)
	expireTime := time.Duration(48 * time.Hour)
	signIn := aauth.AngularSignIn(sessionColl, findUser(app), aauth.NewSha512Password, expireTime)
	signedIn := aauth.AngularAuth(app.DB(), SessionColl)

	newUser := sj.NewAppHandler(sj.NewUserHandler, app)
	newSeries := sj.NewAppHandler(sj.NewSeriesHandler, app)
	removeSeries := sj.NewAppHandler(sj.RemoveSeriesHandler, app)

	readSeriesOfUser := sj.NewAppHandler(sj.ReadSeriesOfUserHandler, app)

	host := app.Specs.Host
	port := app.Specs.Port
	srvRes := host + ":" + strconv.Itoa(port)

	publicDir := app.Specs.PublicDir
	htmlDir := path.Join(publicDir, "html")

	srv := gin.Default()
	srv.Use(sj.Serve("/", sj.LocalFile(htmlDir, false)))
	srv.Static("/public", publicDir)
	srv.GET("/SignIn/:name/:pass", signIn)
	srv.POST("/User", newUser)
	srv.POST("/Series", signedIn, newSeries)
	srv.DELETE("/Series/:id", signedIn, removeSeries)
	srv.GET("/SeriesOfUser", signedIn, readSeriesOfUser)

	srv.Run(srvRes)
}

func findUser(app sj.AppContext) aauth.FindUser {
	return func(name string) (aauth.User, error) {
		return sj.FindUser(app.DB(), name)
	}
}
