package sql

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"images/logutil"
	"io"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB1 *sql.DB

func init_token() {
	logutil.Info("token数据库开始链接")
	var err error
	DB1, err = sql.Open("sqlite3", "./sql/db.token")
	if err != nil {
		logutil.Fatal("token数据库链接失败: %v", err)
	}
	_, err = DB1.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		logutil.Fatal("开启外键约束失败:%v", err)
	}
	// 检查连接是否成功
	err = DB1.Ping()
	if err != nil {
		logutil.Fatal("token数据库链接失败: %v", err)
	}
	// 确定表存在吗
	token_name := "token"
	user_name := "user"
	exists_token, err_token := tableExists(DB1, token_name)
	exists_user, err_user := tableExists(DB1, user_name)
	if err_user != nil {
		logutil.Fatal("token数据库确定user表报错: %v", err_token)
	}
	if !exists_user {
		createQuery := `CREATE TABLE user (id INTEGER PRIMARY KEY AUTOINCREMENT,username TEXT NOT NULL UNIQUE,password TEXT NOT NULL,userid TEXT NOT NULL UNIQUE);`
		_, err := DB1.Exec(createQuery)
		if err != nil {
			logutil.Fatal("token数据库创建失败: %v", err)
		} else {
			logutil.Info("token数据库user表创建成功")
		}
	}
	if err_token != nil {
		logutil.Fatal("token数据库确定token表报错: %v", err_token)
	}
	if !exists_token {
		createQuery := `CREATE TABLE token (id INTEGER PRIMARY KEY AUTOINCREMENT,userid TEXT NOT NULL,token TEXT NOT NULL,refresh TEXT NOT NULL,timestamp INTEGER NOT NULL,FOREIGN KEY (userid) REFERENCES user(userid));`
		_, err := DB1.Exec(createQuery)
		if err != nil {
			logutil.Fatal("token数据库token表创建失败: %v", err)
		} else {
			logutil.Info("token数据库token表创建失败")
		}
	}

	logutil.Info("token数据库链接成功")
}

func close_token() {
	logutil.Info("开始关闭token数据库链接")
	err := DB1.Close()
	if err != nil {
		logutil.Fatal("token数据库关闭失败")
	}
	logutil.Info("token数据库关闭成功")
}

type USER struct {
	ID       int
	USERNAME string
	PASSWORD string
	USERID   string
}

type TOKEN struct {
	ID        int
	USERID    string
	TOKEN     string
	REFRESH   string
	TIMESTAMP int64
}

// 随机生成USERID
func GenerateUSERID(USERNAME string) string {
	// 获取当前时间戳
	timestamp := time.Now().UnixNano()

	// 生成10位随机数字
	// 生成10位随机数字
	randomBytes := make([]byte, 5) // 5 bytes gives us 10 digits in base 10
	if _, err := rand.Read(randomBytes); err != nil {

		logutil.Error("GenerateUSERID 随机数生成失败")
	}
	randomNum := fmt.Sprintf("%010d", int64(binary.LittleEndian.Uint64(append(randomBytes, 0, 0, 0, 0))&0x7FFFFFFFFFFFFFFF))

	// 组合所有部分并进行哈希处理
	input := fmt.Sprintf("%s%d%s", USERNAME, timestamp, randomNum)
	hash := sha256.Sum256([]byte(input))
	hashString := hex.EncodeToString(hash[:])

	// 取出前40个字符作为fileid
	USERID := hashString[:10]
	return USERID
}

// 加密password
func HashPasswordWithSHA256(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	} else {
		passwordWithSalt := password + salt
		hash := sha256.Sum256([]byte(passwordWithSalt))
		return salt + hex.EncodeToString(hash[:]), nil
	}
}
func generateSalt() (string, error) {
	const saltLength = 16
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		logutil.Error("生成盐值失败:%v", err)
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// VerifyPasswordWithSHA256 验证加盐后的密码是否正确
func VerifyPasswordWithSHA256(hashedPasswordWithSalt, password string) bool {
	salt := hashedPasswordWithSalt[0:32]
	passwordWithSalt := password + salt
	hash := sha256.Sum256([]byte(passwordWithSalt))
	return hashedPasswordWithSalt[32:] == hex.EncodeToString(hash[:])
}

// 检查UERID和USERNAME唯一性
func CheckUserIdOrUsername(content string, typeUSER string) (bool, error) {
	// 先检查USERID是否已存在
	var count int

	if typeUSER == "username" {
		err := DB1.QueryRow("SELECT COUNT(*) FROM user WHERE username =?", content).Scan(&count)
		if err != nil {
			logutil.Error("查询USERID是否存在时出错: %v", err)
			return false, err
		}
	} else if typeUSER == "userid" {
		err := DB1.QueryRow("SELECT COUNT(*) FROM user WHERE userid =?", content).Scan(&count)
		if err != nil {
			logutil.Error("查询USERID是否存在时出错: %v", err)
			return false, err
		}
	} else {

		return false, fmt.Errorf("传入typeUSER参数错误：%v", typeUSER)
	}

	if count > 0 {

		return true, nil
	}
	return false, nil
}

// 查询token
func SelectTokenByUserId(UserId string) (TOKEN, error) {
	sqlString := `SELECT id,userid,token,refresh,timestamp  FROM token where userid = ?;`
	var id int
	var timestamp int64
	var userid, token, refresh string
	row := DB1.QueryRow(sqlString, UserId)
	err := row.Scan(&id, &userid, &token, &refresh, &timestamp)
	if err != nil {
		logutil.Error("SelectTokenByUserId(%v) scan报错：%v", UserId, err)
		return TOKEN{}, err
	}
	var Token TOKEN
	Token.ID = id
	Token.TOKEN = token
	Token.REFRESH = refresh
	Token.USERID = userid
	Token.TIMESTAMP = timestamp
	return Token, nil
}

// 查询user
func SelectUserByUserName(UserName string) (USER, error) {
	sqlString := `SELECT id,username,password,userid  FROM user where username = ?;`
	var id int
	var username, password, userid string
	row := DB1.QueryRow(sqlString, UserName)
	err := row.Scan(&id, &username, &password, &userid)
	if err != nil {
		logutil.Error("SelectUserByUserName(%v) scan报错：%v", UserName, err)
		return USER{}, err
	}
	var user USER
	user.ID = id
	user.PASSWORD = password
	user.USERNAME = username
	user.USERID = userid
	return user, nil
}

// 更新token
func UpdataToken(token TOKEN) (bool, error) {
	tx, err := DB1.Begin()
	if err != nil {
		logutil.Error("UpdataToken 事务开始失败 %v", err)
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // 如果有错误发生，回滚事务
		} else {
			err = tx.Commit() // 没有错误则提交事务
		}
	}()
	query := `UPDATE token SET token = ?,refresh = ?,timestamp=? WHERE userid = ?`
	_, err1 := tx.Exec(query, token.TOKEN, token.REFRESH, token.TIMESTAMP, token.USERID)
	if err1 != nil {
		logutil.Error("UpdataToken语句报错 %v", err1)
		return false, err1
	}
	return true, nil
}

// 插入token
func InsertToken(token TOKEN) (bool, error) {
	tx, err := DB1.Begin()
	if err != nil {
		logutil.Error("InsertToken事务开始失败 %v", err)
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // 如果有错误发生，回滚事务
		} else {
			err = tx.Commit() // 没有错误则提交事务
		}
	}()
	query := `INSERT INTO token (userid,token,refresh,timestamp) VALUES (?, ?, ?,?);`
	_, err1 := tx.Exec(query, token.USERID, token.TOKEN, token.REFRESH, token.TIMESTAMP)
	if err1 != nil {
		logutil.Error("InsertToken插入失败 %v", err1)
		return false, err1
	}
	return true, nil
}

// 插入用户
func InsertUser(user USER) (bool, error) {
	tx, err := DB1.Begin()
	if err != nil {
		logutil.Error("InsertUser事务开始失败 %v", err)
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // 如果有错误发生，回滚事务
		} else {
			err = tx.Commit() // 没有错误则提交事务
		}
	}()
	query := `INSERT INTO user (username,password,userid) VALUES (?, ?, ?);`
	_, err1 := tx.Exec(query, user.USERNAME, user.PASSWORD, user.USERID)
	if err1 != nil {
		logutil.Error("InsertUser插入失败 %v", err1)
		return false, err1
	}
	return true, nil
}

// 生成token
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	tokenRandomLength = 5
)

// 生成token
func GenerateToken(userID string) (string, error) {
	// 获取当前时间戳
	timestamp := time.Now().Unix()
	// 生成随机数据
	randomData, err := generateRandomData(tokenRandomLength)
	if err != nil {
		logutil.Error("generateRandomData报错%v", err)
		return "", err
	}
	// 拼接数据
	data := fmt.Sprintf("%s:%d:%s", userID, timestamp, randomData)
	// 密钥，需保密，用于后续的签名验证
	secretKey := []byte("your_secret_key")
	// 使用HMAC-SHA256算法进行签名
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(data))
	signature := fmt.Sprintf("%x", mac.Sum(nil))
	// 将数据和签名拼接，中间用冒号隔开，作为最终的令牌
	token := fmt.Sprintf("%s:%s", data, signature)
	return token, err
}

// 生成refreshtoken
func GenerateRefreshToken(userID string) (string, error) {
	// 获取当前时间戳
	timestamp := time.Now().Unix()
	// 生成随机数据
	randomData, err := generateRandomData(tokenRandomLength)
	if err != nil {
		logutil.Error("generateRandomData报错%v", err)
		return "", err
	}
	// 拼接数据
	data := fmt.Sprintf("%s:%d:%s:%s", userID, timestamp, randomData, "refresh")
	// 密钥，需保密，用于后续的签名验证
	secretKey := []byte("your_secret_key")
	// 使用HMAC-SHA256算法进行签名
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(data))
	signature := fmt.Sprintf("%x", mac.Sum(nil))
	// 将数据和签名拼接，中间用冒号隔开，作为最终的令牌
	token := fmt.Sprintf("%s:%s", data, signature)
	return token, nil
}

func generateRandomData(length int) (string, error) {
	// 创建一个字节切片用于存储随机生成的数据
	b := make([]byte, length)
	// 使用crypto/rand.Read来填充字节切片，生成加密安全的随机字节
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		// 如果出现错误，可以根据实际情况处理，这里简单返回空字符串，实际应用中可能需要更好的错误处理逻辑
		return "", err
	}
	// 将字节切片中的每个字节转换为可打印的字符（从letterBytes中选取对应范围的字符）
	for i := range b {
		b[i] = letterBytes[int(b[i])%len(letterBytes)]
	}
	return string(b), nil
}
