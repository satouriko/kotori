package main

import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/BurntSushi/toml"
	"strconv"
)

var mux = http.NewServeMux()

func main() {

	_, err := toml.DecodeFile("config.toml", &GlobCfg)
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/api", Pong)

	// Example of using a http.FileServer if you want "server-like" rather than "middleware" behavior
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(mux)

	http.ListenAndServe(":" + strconv.FormatInt(GlobCfg.PORT, 10), n)
}
