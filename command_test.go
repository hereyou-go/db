package db

import (
	mdb "database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hereyou-go/logs"
)

func TestBuild(t *testing.T) {
	cmd := NewCommand("SELECT * FROM `tablename` WHERE `colname` BETWEEN :a AND :b")
	cmd.SetParam("a", nil)
	cmd.SetParam("b", nil)
	sql, params, err := cmd.Build()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n[Build SQL]:%s\n[Build Params]:%+v", sql, params)
}

func TestInsertBuild(t *testing.T) {
	cmd := Insert("tablename")
	cmd.SetParam("col", "val")
	sql, params, err := cmd.Build()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n[Build SQL]:%s\n[Build Params]:%+v", sql, params)
}

func TestUpdateBuild(t *testing.T) {
	cmd := Update("tablename", "id")
	cmd.SetParam("col", "val")
	cmd.SetParam("id", "id")
	sql, params, err := cmd.Build()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n[Build SQL]:%s\n[Build Params]:%+v", sql, params)
}

func TestQuery(t *testing.T) {
	sql := "SELECT `file_id` FROM `filess`"
	db, err := mdb.Open("mysql", "root:rootroot123@tcp(noteapp.db:3306)/noteapp?autocommit=false")
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		logs.Fatal(err)
		return
	}
	defer db.Close()
	cmd := NewCommand(sql)
	rows, err := cmd.Query(db)
	if err != nil {
		t.Fatal(err)
	}
	var id int64
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
			return
		}
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}

}
