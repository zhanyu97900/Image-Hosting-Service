//go:build js && wasm
// +build js,wasm

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"syscall/js"
)

// Sha256Encrypt 对输入的字符串进行SHA-256加密并返回十六进制字符串结果
func Sha256Encrypt(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return ""
	}
	input := args[0].String()
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func main() {
	c := make(chan struct{}, 0)
	// 将Go函数注册为JavaScript可调用的函数
	js.Global().Set("sha256Encrypt", js.FuncOf(Sha256Encrypt))
	<-c
}
