package db

import "github.com/hereyou-go/logs"

type UpdateCommand interface {
	Command
	Executable
}

type DBUpdateCommand struct {
	*DBCommand
	table string
	keys  []string
}

func Update(table string, keys ...string) UpdateCommand {
	return &DBUpdateCommand{
		DBCommand: NewCommand(),
		table:     table,
		keys:      keys,
	}
}

func (cmd *DBUpdateCommand) Build() (sql string, params []interface{}, err error) {
	if cmd.keys == nil || len(cmd.keys) == 0 {
		return "", nil, logs.NewError("", "key(s) is required, but is null.")
	}
	fields := ""
	condition := ""
	foundKeys := false
	for col := range cmd.parameters {
		isKey := false
		for _, key := range cmd.keys {
			if key == col {
				foundKeys = true
				isKey = true
				break
			}
		}
		if isKey {
			continue
		}
		if fields != "" {
			fields += ", "
		}
		fields += "`" + col + "`=:" + col
	}
	if !foundKeys {
		return "", nil, logs.NewError("", "key(s) is set, but not found in parameters.")
	}
	for _, col := range cmd.keys {
		if condition != "" {
			condition += "AND "
		}
		condition += "`" + col + "`=:" + col
	}
	cmd.template = "UPDATE `" + cmd.table + "` SET " + fields + " WHERES " + condition + ""
	return cmd.DBCommand.Build()
}
