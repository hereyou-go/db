package db

import "testing"

func TestBuild(t *testing.T) {
	cmd := NewCommand(nil, "SELECT * FROM `tablename` WHERE `colname` BETWEEN :a AND :b")
	cmd.SetParam("a", nil)
	cmd.SetParam("b", nil)
	sql, _, err := cmd.Build()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Build SQL] " + sql)
}
