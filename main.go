package kotori

import (
	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
	"time"
)

var mux = httprouter.New()
var db *gorm.DB
var globalSessions *session.Manager
var startTime time.Time

func main() {

	startTime = time.Now()

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

	mux.GET("/v2", Pong)
	mux.GET("/v2/status", Status)
	mux.GET("/v2/comment", ListComment)
	mux.POST("/v2/comment", CreateComment)
	mux.DELETE("/v2/comment/:id", DeleteComment)
	mux.POST("/v2/auth", Login)
	mux.DELETE("/v2/auth", Logout)
	mux.PUT("/v2/user/:id", EditUserSetHonor)
	mux.GET("/v2/index", ListIndex)
	mux.GET("/v2/index/:id", GetIndex)
	mux.POST("/v2/index", CreateIndex)
	mux.PUT("/v2/index/:id", EditIndex)
	mux.DELETE("/v2/index/:id", DeleteIndex)
	mux.GET("/v2/post", ListPost)
	mux.GET("/v2/post/:id", GetPost)
	mux.POST("/v2/post", CreatePost)
	mux.PUT("/v2/post/:id", EditPost)
	mux.DELETE("/v2/post/:id", DeletePost)

	c := cors.New(cors.Options{
		AllowedOrigins:   GlobCfg.ALLOW_ORIGIN,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Query-By"},
	})
	handler := c.Handler(mux)

	n := negroni.New()
	n.UseHandler(handler)

	http.ListenAndServe(":"+strconv.FormatInt(GlobCfg.PORT, 10), n)
}
