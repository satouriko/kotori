package main

import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/BurntSushi/toml"
	"strconv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/julienschmidt/httprouter"
)

var mux = httprouter.New()
var db *gorm.DB

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

	mux.GET("/api", Pong)
	mux.GET("/api/comment", GetComment)
	mux.POST("/api/comment", StoreComment)
	mux.DELETE("/api/comment/:id", DeleteComment)
	mux.POST("/api/auth", Login)
	mux.DELETE("/api/auth", Logout)
	mux.GET("/api/index", GetIndex)
	mux.POST("/api/index", StoreIndex)
	mux.PUT("/api/index/:id", UpdateIndex)
	mux.DELETE("/api/index/:id", DeleteIndex)
	mux.GET("/api/post", GetPost)
	mux.POST("/api/post", StorePost)
	mux.PUT("/api/post/:id", UpdatePost)
	mux.DELETE("/api/post/:id", DeletePost)

	mux.ServeFiles("/static/*filepath", http.Dir("static"))

	// Example of using a http.FileServer if you want "server-like" rather than "middleware" behavior
	//mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(mux)

	http.ListenAndServe(":"+strconv.FormatInt(GlobCfg.PORT, 10), n)
}
