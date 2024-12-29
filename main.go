package main

import (
	"context"
	"images/logutil"
	"images/server"
	"images/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// println("欢迎使用sqliteview服务")
	// println("请访问localhost:5000,或者本机IP:5000")
	// 创建默认的路由引擎
	logutil.Init()
	sql.Init()

	server.Init()
	go func() {
		if err := server.SERVER.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logutil.Info("监听: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logutil.Info("关闭服务中")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.SERVER.Shutdown(ctx); err != nil {
		logutil.Info("服务关闭: %s\n", err)
	}

	logutil.Info("成功关闭服务")
	defer sql.Close()
	defer logutil.Close()

}
