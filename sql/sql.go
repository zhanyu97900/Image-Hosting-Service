package sql

import (
	"database/sql"
	"errors"
	"images/logutil"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	logutil.Info("数据库开始链接")
	var err error
	DB, err = sql.Open("sqlite3", "./sql/db.images")
	if err != nil {
		logutil.Fatal("数据库链接失败: %v", err)
	}

	// 检查连接是否成功
	err = DB.Ping()
	if err != nil {
		logutil.Fatal("数据库链接失败: %v", err)
	}
	logutil.Info("数据库链接成功")
}

func Close() {
	logutil.Info("开始关闭数据库链接")
	err := DB.Close()
	if err != nil {
		logutil.Fatal("数据库关闭失败")
	}
	logutil.Info("数据库关闭成功")
}

type Image struct {
	Id      int    `json:"-"`
	FileId  string `json:"FileId"`
	Sha256  string `json:"-"`
	Name    string `json:"Name"`
	Path    string `json:"-"`
	Removed int    `json:"-"`
}

// 查询Images
func Search_sql(fileid string) ([]Image, error) {
	var Image_lst []Image
	var sql_str string
	if fileid == "all" {
		sql_str = `SELECT id,fileid,sha256,name,path,removed FROM images where removed = ?;`
		rows, err := DB.Query(sql_str, 0)
		if err != nil {
			logutil.Error("sql: search(%v) Query报错：%v", fileid, err)
			return nil, err
		}
		defer rows.Close()
		var id int
		var fileid, sha256, name, path string
		var removed int
		for rows.Next() {
			err := rows.Scan(&id, &fileid, &sha256, &name, &path, &removed)
			if err != nil {
				logutil.Error("sql: search(%v) scan报错：%v", fileid, err)
				return nil, err
			}
			var img Image
			img.Id = id
			img.FileId = fileid
			img.Name = name
			img.Path = path
			img.Sha256 = sha256
			img.Removed = removed

			Image_lst = append(Image_lst, img)

		}

	} else {
		sql_str = `SELECT id,fileid,sha256,name,path,removed FROM images WHERE fileid=? AND removed = ?;`
		row := DB.QueryRow(sql_str, fileid, 0)
		var id int
		var fileid, sha256, name, path string
		var removed int
		err := row.Scan(&id, &fileid, &sha256, &name, &path, &removed)
		if err != nil {
			logutil.Error("sql: search(%v) scan报错：%v", fileid, err)
			return nil, err
		}
		var img Image
		img.Id = id
		img.FileId = fileid
		img.Name = name
		img.Path = path
		img.Sha256 = sha256
		img.Removed = removed
		Image_lst = append(Image_lst, img)

	}
	return Image_lst, nil
}

// searchBysha256
func Search_sha256(sha256_ string) ([]Image, error) {
	var images []Image
	sql_str := `SELECT id,fileid,sha256,name,path,removed FROM images where sha256 = ? AND removed=?;`
	rows, err := DB.Query(sql_str, sha256_, 0)
	if err != nil {
		logutil.Error("sql: search_sha256(%v) Query报错：%v", sha256_, err)
		return nil, err
	}
	defer rows.Close()
	var id int
	var fileid, sha256, name, path string
	var removed int
	for rows.Next() {
		err := rows.Scan(&id, &fileid, &sha256, &name, &path, &removed)
		if err != nil {
			logutil.Error("sql: search_sha256(%v) scan报错：%v", sha256_, err)
			return nil, err
		}
		var img Image
		img.Id = id
		img.FileId = fileid
		img.Name = name
		img.Path = path
		img.Sha256 = sha256
		img.Removed = removed

		images = append(images, img)

	}
	return images, nil

}

// 删除方法
func Remove_sql(fileid string) (bool, error) {
	images, err := Search_sql(fileid)
	if err != nil {
		logutil.Error("sql remove_sql: 删除失败-目标不存在 %v", err)
		return false, err
	}
	if len(images) == 0 {
		err2 := errors.New("目标不存在")
		logutil.Error("sql remove_sql: 尝试删除不存在的目标 %s", fileid)
		return false, err2
	}
	sql_str := `UPDATE images SET removed = ? WHERE fileid = ?`
	_, err_1 := DB.Exec(sql_str, 1, fileid)
	if err_1 != nil {
		logutil.Error("sql remove_sql: 删除失败 %v", err_1)
		return false, err_1
	}
	return true, nil

}

// 查询现在被删除的内容
func removed_sql() ([]Image, error) {
	var images []Image
	sql_str := `SELECT id,fileid,sha256,name,path,removed FROM images WHERE removed = ?;`
	rows, err := DB.Query(sql_str, 1)
	if err != nil {
		logutil.Error("sql: removed_sql Query报错：%v", err)
		return nil, err

	}
	defer rows.Close()
	var id int
	var fileid, sha256, name, path string
	var removed int
	for rows.Next() {
		var img Image

		err := rows.Scan(&id, &fileid, &sha256, &name, &path, &removed)
		if err != nil {
			logutil.Error("sql: removed_sql scan报错：%v", err)
			return nil, err
		}
		img.FileId = fileid
		img.Id = id
		img.Name = name
		img.Path = path
		img.Removed = removed
		img.Sha256 = sha256

		images = append(images, img)
	}
	return images, nil
}

// 插入函数
func Insert_sql(images Image) (bool, error) {
	imgs, err2 := removed_sql()
	if err2 != nil {
		logutil.Error("sql insert_sql 报错 %v", err2)
		return false, err2
	}
	tx, err := DB.Begin()
	if err != nil {
		logutil.Error("sql insert_sql 事务开始失败 %v", err)
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // 如果有错误发生，回滚事务
		} else {
			err = tx.Commit() // 没有错误则提交事务
		}
	}()
	if len(imgs) != 0 {
		img := imgs[0]
		sql_str := `UPDATE images SET fileid=?,name=?,path=?,sha256=?,removed = ? WHERE id=? AND fileid = ?`
		_, err := tx.Exec(sql_str, images.FileId, images.Name, images.Path, images.Sha256, 0, img.Id, img.FileId)
		if err != nil {
			logutil.Error("sql Insert_sql updata 语句报错%v", err)
			return false, err
		}
	} else {
		query := `INSERT INTO images (fileid,name,path,sha256,removed) VALUES (?, ?, ?, ?,?)`
		_, err := tx.Exec(query, images.FileId, images.Name, images.Path, images.Sha256, images.Removed)
		if err != nil {
			logutil.Error("sql insert_sql insert语句报错 %v", err)
			return false, err
		}
	}

	return true, nil
}

// prepare表
type Prepare struct {
	Id        int
	Fileid    string
	Timestamp int64
	Name      string
	Upload    int
}

// prepare 查询fileid
func Sql_prepare_fileid(fileid_ string) ([]Prepare, error) {
	var prepares []Prepare
	sql_str := `SELECT id,fileid,timestamp,name,upload FROM prepare WHERE fileid = ?;`
	rows, err := DB.Query(sql_str, fileid_)
	if err != nil {
		logutil.Error("sql: Sql_prepare_fileid 报错：%v", err)
		return nil, err

	}
	defer rows.Close()
	var id, upload int
	var timestamp int64
	var name, fileid string
	for rows.Next() {
		var prepare Prepare

		err := rows.Scan(&id, &fileid, &timestamp, &name, &upload)
		if err != nil {
			logutil.Error("sql: removed_sql scan报错：%v", err)
			return nil, err
		}

		prepare.Id = id
		prepare.Fileid = fileid
		prepare.Upload = upload
		prepare.Name = name
		prepare.Timestamp = timestamp

		prepares = append(prepares, prepare)
	}
	return prepares, nil
}

// prepare插入fileid
func Sql_prepare_add_fileid(fileid string, name string, timestamp_ int64) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		logutil.Error("sql Sql_prepare_add_fileid 事务开始失败 %v", err)
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // 如果有错误发生，回滚事务
		} else {
			err = tx.Commit() // 没有错误则提交事务
		}
	}()

	query := `INSERT INTO prepare (fileid,name,timestamp,upload) VALUES (?, ?, ?, ?);`
	_, err1 := tx.Exec(query, fileid, name, timestamp_, 0)
	if err1 != nil {
		logutil.Error("sql Sql_prepare_add_fileid insert语句报错 %v", err1)
		return false, err1
	}
	return true, nil
}

// prepare 更新upload状态
func Sql_prepare_upload_fileid(fileid_ string, id_ int) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		logutil.Error("sql Sql_prepare_upload_fileid 事务开始失败 %v", err)
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // 如果有错误发生，回滚事务
		} else {
			err = tx.Commit() // 没有错误则提交事务
		}
	}()
	sql_str := `UPDATE prepare SET upload=? WHERE id=? AND fileid = ?`
	_, err1 := tx.Exec(sql_str, 1, id_, fileid_)
	if err1 != nil {
		logutil.Error("sql Insert_sql updata 语句报错%v", err1)
		return false, err1
	}
	return true, nil
}
