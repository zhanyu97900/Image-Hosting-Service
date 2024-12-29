package server

import (
	"images/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

var R *gin.Engine
var SERVER http.Server

func Init() {
	R = gin.Default()
	SERVER.Addr = ":5000"
	SERVER.Handler = R
	R.NoRoute(handlers.NotFound)
	R.LoadHTMLGlob("templates/*")
	R.Static("/static", "./static")
	// api
	R.POST("/api/upload", handlers.Upload_handler)
	R.POST("/api/prepare", handlers.Prepare_handler)
	R.POST("/api/remove", handlers.Remove_handler)
	R.GET("/file/*fileid", handlers.Respone_Images)
	R.GET("/api/allfile", handlers.Respone_all)
	// static
	R.GET("/upload", handlers.Upload)
	R.GET("/admin", handlers.Admin)
	R.GET("/", handlers.Index)
}
