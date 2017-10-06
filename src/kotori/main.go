package main

import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/BurntSushi/toml"
	"strconv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/julienschmidt/httprouter"
	"github.com/astaxie/beego/session"
	"github.com/rs/cors"
)

var mux = httprouter.New()
var db *gorm.DB
var globalSessions *session.Manager

func main() {

	_, err := toml.DecodeFile("config.toml", &GlobCfg)
	if err != nil {
		panic(err)
	}

	db, err = gorm.Open("sqlite3", "core.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Index{}, &User{}, &Comment{}, &Post{})

	globalSessions, _ = session.NewManager("memory", &session.ManagerConfig{CookieName: "kotoriCoreSession", EnableSetCookie: true, Gclifetime: 3600})
	go globalSessions.GC()

	mux.GET("/api", Pong)
	mux.GET("/api/comment", ListComment)
	mux.POST("/api/comment", CreateComment)
	mux.DELETE("/api/comment/:id", DeleteComment)
	mux.POST("/api/auth", Login)
	mux.DELETE("/api/auth", Logout)
	mux.PUT("/api/user/:id", EditUserSetHonor)
	mux.GET("/api/index", ListIndex)
	mux.POST("/api/index", CreateIndex)
	mux.PUT("/api/index/:id", EditIndex)
	mux.DELETE("/api/index/:id", DeleteIndex)
	mux.GET("/api/post", ListPost)
	mux.GET("/api/post/:id", GetPost)
	mux.POST("/api/post", CreatePost)
	mux.PUT("/api/post/:id", EditPost)
	mux.DELETE("/api/post/:id", DeletePost)

	mux.ServeFiles("/static/*filepath", http.Dir("static"))

	c := cors.New(cors.Options{
		AllowedOrigins: GlobCfg.ALLOW_ORIGIN,
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(handler)

	http.ListenAndServe(":"+strconv.FormatInt(GlobCfg.PORT, 10), n)
}
