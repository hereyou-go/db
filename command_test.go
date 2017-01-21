package db

import "testing"

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
