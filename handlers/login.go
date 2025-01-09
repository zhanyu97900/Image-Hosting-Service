package handlers

import (
	"images/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// login 登录
func LoginApi(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	exists, err := sql.CheckUserIdOrUsername(username, "username")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "服务器错误无法检查账号是否存在"})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": "账号不存在"})
		return
	}
	user, err2 := sql.SelectUserByUserName(username)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "服务器错误无法查询用户"})
		return
	}
	same := sql.VerifyPasswordWithSHA256(user.PASSWORD, password)
	if !same {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": "密码错误"})
		return
	}
	token, err4 := sql.GenerateToken(user.USERID)
	refresh, err5 := sql.GenerateRefreshToken(user.USERID)
	if err4 != nil && err5 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "token生成失败"})
		return
	}
	var Token sql.TOKEN
	var timestamp int64 = time.Now().Unix()
	Token.USERID = user.USERID
	Token.REFRESH = refresh
	Token.TOKEN = token
	Token.TIMESTAMP = timestamp
	ok, err6 := sql.UpdataToken(Token)
	if err6 != nil || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "token储存失败"})
		return
	}
	c.SetCookie("token", token, 7200, "/", "", false, true)
	c.SetCookie("refresh", refresh, 10800, "/", "", false, true)
	c.SetCookie("sign", user.USERID, 10800, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "ok"})
}

// 注册模块
func RegisterApi(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	exists, err_exists := sql.CheckUserIdOrUsername(username, "username")
	if err_exists != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "服务器错误无法确定用户"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": "用户已存在"})
		return
	}
	var userid string
	for i := 1; i < 5; i++ {
		userid = sql.GenerateUSERID(username)
		exists, err_useid := sql.CheckUserIdOrUsername(userid, "userid")
		if err_useid != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "无法创建userid"})
			return
		}
		if !exists {
			break
		}
	}
	newPassword, err_hash := sql.HashPasswordWithSHA256(password)
	if err_hash != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "服务器错误"})
		return
	}
	var user sql.USER
	var token sql.TOKEN
	user.PASSWORD = newPassword
	user.USERID = userid
	user.USERNAME = username
	token.USERID = userid
	token.REFRESH = ""
	token.TOKEN = ""
	token.TIMESTAMP = int64(0)
	status, err_insert := sql.InsertUser(user)
	status_token, err_token := sql.InsertToken(token)
	if err_insert != nil || err_token != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "服务器错误"})
		return
	}
	if !status || !status_token {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "data": "创建用户失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "用户创建成功"})
}
