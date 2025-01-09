package middleware

import (
	"images/logutil"
	"images/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginCookie(c *gin.Context) {
	ip := c.ClientIP()

	// 获取并检查 cookies
	userid, err_userid := c.Cookie("sign")
	token, err_token := c.Cookie("token")
	refresh, err_refresh := c.Cookie("refresh")

	if err_userid != nil || err_token != nil || err_refresh != nil {
		logutil.Error("%v 获取cookie报错 userid:%v  token:%v  refresh:%v", ip, err_userid, err_token, err_refresh)
		handleError(c, "/login")
		return
	}

	if userid == "" || token == "" || refresh == "" {
		logutil.Info("%v 获取cookie为空 userid:%v  token:%v  refresh:%v", ip, userid, token, refresh)
		handleError(c, "/login")
		return
	}

	// 验证 token
	Token, err := sql.SelectTokenByUserId(userid)
	timestamp := time.Now().Unix()
	if err != nil {
		handleError(c, "/login")
		return
	}

	if token != Token.TOKEN {
		logutil.Info("%v 获取token:%v 与token不同:%v", ip, token, Token.TOKEN)
		handleError(c, "/login")
		return
	}

	// 检查 token 是否过期
	if Token.TIMESTAMP+int64(7200) < timestamp {
		// 如果 refresh token 还有效，则更新 token
		if Token.TIMESTAMP+int64(10800) >= timestamp && refresh == Token.REFRESH {
			updateTokens(c, userid, ip)
		} else {
			logutil.Info("%v 无法刷新，刷新token过期", ip)
			handleError(c, "/login")
		}
		return
	}

	// 如果一切都正常，继续执行后续处理器
	c.Next()
}

func handleError(c *gin.Context, target string) {
	if target == "500.html" {
		c.HTML(http.StatusInternalServerError, target, nil)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, target)
	}
	c.Abort()
}

func updateTokens(c *gin.Context, userid string, ip string) {
	tokenNew, err_token := sql.GenerateToken(userid)
	refreshNew, err_refresh := sql.GenerateRefreshToken(userid)

	if err_refresh != nil || err_token != nil {
		logutil.Error("%v token更新过程失败 tokenNew:%v refreshNew:%v", ip, err_token, err_refresh)
		handleError(c, "500.html")
		return
	}

	var TokenNew sql.TOKEN
	TokenNew.TIMESTAMP = time.Now().Unix()
	TokenNew.REFRESH = refreshNew
	TokenNew.TOKEN = tokenNew
	TokenNew.USERID = userid

	status_token, err := sql.UpdataToken(TokenNew)
	if err != nil || !status_token {
		logutil.Error("%v token更新失败: %v", ip, err)
		handleError(c, "/login")
		return
	}

	// 设置新的 cookies
	c.SetCookie("token", tokenNew, 7200, "/", "", false, true)
	c.SetCookie("refresh", refreshNew, 10800, "/", "", false, true)
	c.SetCookie("sign", userid, 10800, "/", "", false, true)

	// 继续执行后续处理器
	c.Next()
}

func CookieMiddleware() gin.HandlerFunc {
	return LoginCookie
}
