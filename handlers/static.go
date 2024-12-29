package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
func Upload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
}

func Admin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}
