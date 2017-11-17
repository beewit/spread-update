package test

import (
	"database/sql"
	"github.com/beewit/beekit/sqlite"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread-update/update"
	"github.com/pkg/errors"
	"testing"
)

func TestUpdateDB(t *testing.T) {
	err := CheckUpdate()
	if err != nil {
		t.Error(err)
	}
	println("-------------结束UpdateDB------------")
}

//func TestInsertDB(t *testing.T) {
//	err := InitNewDB()
//	if err != nil {
//		t.Error(err.Error())
//		return
//	}
//	defer NewDB.Close()
//	for i := 0; i < 100; i++ {
//		m := map[string]interface{}{}
//		m["id"] = i
//		m["iden"] = "123456"
//		m["nick_name"] = "123456"
//
//		NewDB.InsertMap("account_http_cache", m)
//	}
//
//}

func GetVersion() (int, error) {
	err := InitNewDB()
	if err != nil {
		return 0, err
	}
	defer NewDB.Close()
	sql := `SELECT version FROM local_version WHERE id=1 LIMIT 1`
	m, err := NewDB.Query(sql)
	if err != nil {
		return 0, err
	}
	if len(m) <= 0 {
		return 0, errors.New("数据库版本错误")
	}
	return convert.MustInt(m[0]["version"]), nil
}

func UpdateVersion(version int) error {
	sql := `DELETE FROM local_version WHERE id=1;INSERT INTO local_version(id,version) VALUES(1,?)`
	_, err := NewDB.Update(sql, version)
	if err != nil {
		return err
	}
	return err
}

//查询所有表
func SelectTables() ([]map[string]interface{}, error) {
	sql := "select * from sqlite_master WHERE type = 'table'"
	return OldDB.Query(sql)
}

//查询所有表所有字段
func SelectFail(table string) ([]map[string]interface{}, error) {
	sql := "PRAGMA table_info(?)"
	return OldDB.Query(sql, table)
}

func QueryOldDBData(pageIndex int, table string) (*utils.PageData, error) {
	return OldDB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     table,
		PageIndex: pageIndex,
		PageSize:  100,
	})
}

//插入旧版数据
func InsertOldDBData(pageIndex int, table string) error {
	page, err := QueryOldDBData(pageIndex, table)
	if err != nil {
		return err
	}
	if page.Count <= 0 {
		return nil
	}
	//执行数据转移
	for _, v := range page.Data {
		_, err = NewDB.InsertMap(table, v)
		if err != nil {
			return err
		}
	}
	if page.PageIndex < page.PageSize {
		pageIndex++
		err = InsertOldDBData(pageIndex, table)
		if err != nil {
			return err
		}
	}
	return nil
}

//查询表所有数据并导入新表数据
func ImportData() error {
	//1、查询旧版数据库所有数据库表，进行循环操作获取表字段
	m, err := SelectTables()
	if err != nil {
		return err
	}
	if len(m) < 0 {
		return errors.New("无表结构")
	}
	for _, v := range m {
		tableName := convert.ToObjStr(v["name"])
		err = InsertOldDBData(1, tableName)
		if err != nil {
			println(err)
			continue
		}
	}
	return nil
}

func CheckUpdate() (err error) {
	var version int
	version, err = GetVersion()
	if err != nil {
		return
	}
	_, err = update.DBUpdate(update.Version{Major: 0, Minor: 0, Patch: version}, func(fileNames []string, rel update.Release) {
		if len(fileNames) > 0 {
			for _, name := range fileNames {
				if name == "spread.db" {
					err2 := InitNewDB()
					if err2 != nil {
						println("InitNewDB ERROR", err.Error())
						return
					}
					defer NewDB.Close()
					err2 = InitOldDB()
					if err2 != nil {
						println("InitOldDB ERROR", err.Error())
						return
					}
					defer OldDB.Close()
					err2 = ImportData()
					if err2 != nil {
						println("InitOldDB ERROR", err.Error())
						return
					}
					UpdateVersion(rel.Patch)
				}
			}
		}
	})
	return
}

var (
	OldDB *sqlite.SqlConnPool
	NewDB *sqlite.SqlConnPool
)

func InitNewDB() (err error) {
	var flog bool
	flog, err = utils.PathExists("spread.db")
	if err != nil {
		return
	}
	if !flog {
		err = errors.New("更新数据库文件已损坏无法更新")
		return
	}
	NewDB = &sqlite.SqlConnPool{
		DriverName:     "sqlite3",
		DataSourceName: "spread.db",
	}
	NewDB.SqlDB, err = sql.Open(NewDB.DriverName, NewDB.DataSourceName)
	if err != nil {
		return
	}
	return
}

/**
特别注意，新版本必须包含兼容老版本数据库结构
*/
func InitOldDB() (err error) {
	var OldFlog bool
	OldFlog, err = utils.PathExists("spread.db.old")
	if err != nil {
		return
	}
	if !OldFlog {
		err = errors.New("更新数据库文件已损坏无法更新")
		return
	}

	OldDB = &sqlite.SqlConnPool{
		DriverName:     "sqlite3",
		DataSourceName: "spread.db.old",
	}
	OldDB.SqlDB, err = sql.Open(OldDB.DriverName, OldDB.DataSourceName)
	if err != nil {
		return
	}
	return
}
