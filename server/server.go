package server

import (
	"images/handlers"
	"images/middleware"
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

	protected := R.Group("/")

	// api
	R.POST("/api/login", handlers.LoginApi)
	R.POST("/api/register", handlers.RegisterApi)
	R.GET("/file/*fileid", handlers.Respone_Images)
	// static
	R.GET("/login", handlers.LoginPage)
	R.GET("/register", handlers.RegisterPage)
	// 设置保护
	protected.Use(middleware.CookieMiddleware())
	{
		// static
		protected.GET("/upload", handlers.Upload)
		protected.GET("/admin", handlers.Admin)
		protected.GET("/", handlers.Index)
		//api
		protected.POST("/api/upload", handlers.Upload_handler)
		protected.POST("/api/prepare", handlers.Prepare_handler)
		protected.POST("/api/remove", handlers.Remove_handler)
		protected.GET("/api/allfile", handlers.Respone_all)
	}

}
