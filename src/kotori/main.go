package main

import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/BurntSushi/toml"
	"strconv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var mux = http.NewServeMux()

func main() {

	_, err := toml.DecodeFile("config.toml", &GlobCfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("sqlite3", "core.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Index{}, &User{}, &Comment{}, &Post{})


	mux.HandleFunc("/api", Pong)

	// Example of using a http.FileServer if you want "server-like" rather than "middleware" behavior
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(mux)

	http.ListenAndServe(":" + strconv.FormatInt(GlobCfg.PORT, 10), n)
}
