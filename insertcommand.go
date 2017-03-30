package db
import (
	"database/sql"

	"github.com/hereyou-go/logs"
)
type InsertCommand interface {
	Command
	Executable
}

type DBInsertCommand struct {
	DBCommand
	table string
}

func Insert(table string) InsertCommand {
	return &DBInsertCommand{
		DBCommand: *NewCommand(),
		table:     table,
	}
}

func (cmd *DBInsertCommand) Build() (sql string, params []interface{}, err error) {
	fields := ""
	vals := ""
	for col := range cmd.DBCommand.parameters {
		if fields != "" {
			fields += ", "
			vals += ", "
		}
		fields += "`" + col + "`"
		vals += ":" + col
	}
	cmd.DBCommand.template = "INSERT INTO `" + cmd.table + "` (" + fields + ") VALUES (" + vals + ")"
	return cmd.DBCommand.Build()
}

func (cmd *DBInsertCommand) Exec(conn Connection) (sql.Result, error) {
	query, args, err := cmd.Build()
	if err != nil {
		return nil, logs.Wrap(err)
	}
	logs.Debug("\n[SQL] %v \n[Parameters] %+v", query, args)
	return conn.Exec(query, args...)
}