package goutils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestGetSchema(t *testing.T) {
	db, err := sql.Open("mysql", "root:112233@tcp(127.0.0.1:3305)/test?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println("数据库链接错误", err)
		return
	}
	// 延迟到函数结束关闭链接
	defer db.Close()

	// ds, err := GetTableSchema(db, "testdb", "test_table_1")
	// ds, err := GetTableSchemas(db, "testdb")
	ds, err := GetColumnSchemas(db, "testdb", "test_table_1")

	if err != nil {
		t.Error(err)
	}

	fmt.Println(ds)
}
