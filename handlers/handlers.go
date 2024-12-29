package handlers

import (
	"images/image"
	"images/logutil"
	"images/sql"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// prepare行
func Prepare_handler(c *gin.Context) {
	filename := c.PostForm("filename")
	sha256 := c.PostForm("sha256")
	fileid := image.GenerateFileID(filename)
	images, err := sql.Search_sha256(sha256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
		})
		return
	}
	time_stamp_ := time.Now().Unix()
	code, err := sql.Sql_prepare_add_fileid(fileid, filename, time_stamp_)
	if err != nil || !code {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
		})
		return
	}
	if len(images) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 2, "data": images[0].FileId, "newFileid": fileid,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0, "data": fileid,
		})
	}
}

// upload文件
func Upload_handler(c *gin.Context) {
	fileid := c.PostForm("fileid")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "检索文件时出错"})
		return
	}
	if fileid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": "fileid缺少"})
		return
	}
	src, err1 := file.Open()
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": "file打开失败"})
		return
	}
	defer src.Close()
	file_path, file_sha256, err_save := image.Save_image(src, file.Filename, fileid)
	if err_save != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "file保存失败"})
		return
	}
	prepare, err_sql_fileid := sql.Sql_prepare_fileid(fileid)
	if err_sql_fileid != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "fileid不存在"})
		return
	}
	if prepare[0].Timestamp+3600 < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": "fileid已超时"})
		return
	}
	var images sql.Image
	images.FileId = fileid
	images.Name = file.Filename
	images.Path = file_path
	images.Removed = 0
	images.Sha256 = file_sha256
	code, err_insert := sql.Insert_sql(images)
	code_p, _ := sql.Sql_prepare_upload_fileid(fileid, prepare[0].Id)
	if err_insert != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "文件保存在数据库失败"})
		return
	}
	if code && code_p {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": fileid})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "文件保存在数据库失败"})
		return
	}
}

// remove文件
func Remove_handler(c *gin.Context) {
	fileid := c.PostForm("fileid")
	images, err := sql.Search_sql(fileid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "删除失败"})
		return
	} else if len(images) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": "file文件不存在"})
		return
	}

	removed, err1 := sql.Remove_sql(images[0].FileId)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "删除失败"})
		return
	}
	if !removed {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "删除失败"})
		return
	}
	if removed {
		removed, err1 := image.Remove_image(images[0].Path)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "删除失败"})
			return
		}
		if !removed {
			logutil.Error("%v 文件未删除", images[0].Name)
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": "删除成功--但是文件未删除"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "删除成功"})
}

// 响应文件
func Respone_Images(c *gin.Context) {
	fileid := c.Param("fileid")
	if len(fileid) > 0 && fileid[0] == '/' {
		fileid = fileid[1:]
	}
	if fileid == "" {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	images, err_file := sql.Search_sql(fileid)
	if err_file != nil {
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}
	if len(images) == 0 {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	} else {
		file, err := os.Open(images[0].Path)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
		defer file.Close()
		ext := image.Get_file_extension(images[0].Path)
		mime_ := image.Get_content_type(ext)
		c.Header("Content-Type", mime_)
		if _, err := io.Copy(c.Writer, file); err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
	}

}

// 响应全部文件
func Respone_all(c *gin.Context) {
	images, err := sql.Search_sql("all")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "错误"})
		return
	}
	if len(images) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "data": images})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": images})
		return
	}
}
