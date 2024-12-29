package logutil

import (
	"log"
	"os"
	"sync"
)

// Logger 是全局的日志记录器
var Logger *log.Logger
var logMutex sync.Mutex // 用于锁定日志写入

// InitializeLogger 初始化日志记录器
func InitializeLogger(logFilePath string) {
	// 打开或创建日志文件
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("打开logfile失败: %v", err)
	}

	// 创建一个新的日志记录器
	Logger = log.New(file, "", log.LstdFlags)
}

// Info 输出信息级别的日志
func Info(format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	Logger.Printf(format, v...)
}

// Error 输出错误级别的日志
func Error(format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	Logger.Printf(format, v...)
}

// Fatal 输出致命错误级别的日志并终止程序
func Fatal(format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	Logger.Fatalf(format, v...)
}

func Init() {
	InitializeLogger("./logutil/log.log")
	Info("日志系统加载完成")
}

func Close() {
	Info("日志模块退出成功")

}
