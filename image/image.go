package image

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"images/logutil"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// 生成fileid
func GenerateFileID(filename string) string {
	// 获取当前时间戳
	timestamp := time.Now().UnixNano()

	// 生成10位随机数字
	// 生成10位随机数字
	randomBytes := make([]byte, 5) // 5 bytes gives us 10 digits in base 10
	if _, err := rand.Read(randomBytes); err != nil {

		logutil.Error("随机数生成失败")
	}
	randomNum := fmt.Sprintf("%010d", int64(binary.LittleEndian.Uint64(append(randomBytes, 0, 0, 0, 0))&0x7FFFFFFFFFFFFFFF))

	// 组合所有部分并进行哈希处理
	input := fmt.Sprintf("%s%d%s", filename, timestamp, randomNum)
	hash := sha256.Sum256([]byte(input))
	hashString := hex.EncodeToString(hash[:])

	// 取出前40个字符作为fileid
	fileID := hashString[:40]
	return fileID
}

// image 保存
func Save_image(file multipart.File, filename string, fileid string) (string, string, error) {
	var basePath_ = "./image/"
	time_now := time.Now().Format("2006-01-02")
	filepath_ := filepath.Join(basePath_, time_now, fmt.Sprintf("%s-%s", fileid, filename))
	dir := filepath.Dir(filepath_)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logutil.Error("创建文件失败: %v", err)
		return "", "", err
	}
	dst, err := os.OpenFile(filepath_, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logutil.Error("Save_image 报错 %v", err)
		return "", "", err
	}
	defer dst.Close()
	hasher := sha256.New()

	// 复制文件内容到目标文件和哈希计算器
	teeReader := io.TeeReader(file, hasher)
	_, err = io.Copy(dst, teeReader)
	if err != nil {
		logutil.Error("Save_image 报错 %v", err)
		return "", "", err
	}
	hashValue := fmt.Sprintf("%x", hasher.Sum(nil))
	return filepath_, hashValue, nil
}

// 删除 image
func Remove_image(file_path string) (bool, error) {
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		logutil.Error("文件 '%s' 不存在。\n", file_path)
		return false, err
	} else if err != nil {
		logutil.Error("检查文件时出错: %v\n", err)
	}
	err := os.Remove(file_path)
	if err != nil {
		logutil.Error("无法删除文件 '%s': %v\n", file_path, err)
		return false, err
	}
	return true, nil
}

// 获取MIME信息
func Get_content_type(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".tiff":
		return "image/tiff"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// 文件格式
func Get_file_extension(path string) string {
	ext := filepath.Ext(path)
	return ext
}
