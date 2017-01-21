package db

type InsertCommand struct {
	*ExecutableCommand
	table string
}

func Insert(table string) *InsertCommand {
	return &InsertCommand{
		ExecutableCommand: newExecutable(),
		table:             table,
	}
}

func (cmd *InsertCommand) Build() (sql string, params []interface{}, err error) {

	fields := ""
	vals := ""
	for col := range cmd.parameters {
		if fields != "" {
			fields += ", "
			vals += ", "
		}
		fields += "`" + col + "`"
		vals += ":" + col
	}
	cmd.template = "INSERT INTO `" + cmd.table + "` (" + fields + ") VALUES (" + vals + ")"
	return cmd.Command.Build()
}
