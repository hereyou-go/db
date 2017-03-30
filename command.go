package db

import (
	"database/sql"
	"regexp"
	"strings"

	"github.com/hereyou-go/logs"
)

type Command interface {
	Set(name string, value interface{}) error
	Build() (sql string, params []interface{}, err error)
}

type Executable interface {
	Exec(conn Connection) (sql.Result, error)
}

type Queryable interface {
	Query(conn Connection) (*sql.Rows, error)
}

type DBCommand struct {
	conn         Connection
	parameters   map[string]interface{}
	template     string
	rebuild      bool // built
	cachedSQL    string
	cachedParams []interface{}
}

func NewCommand(template ...string) *DBCommand {
	cmd := &DBCommand{
		parameters: make(map[string]interface{}),
		rebuild:    true,
	}
	if len(template) > 0 {
		cmd.SetCommand(template[0])
	}
	return cmd
}
func (cmd *DBCommand) SetCommand(template string) error {
	cmd.template = template
	cmd.rebuild = true
	return nil
}

func (cmd *DBCommand) Set(name string, value interface{}) error {
	cmd.parameters[strings.ToLower(name)] = value
	return nil
}

// SetParam is deprecated.
// func (cmd *DBCommand) SetParam(name string, value interface{}) error {
// 	cmd.parameters[strings.ToLower(name)] = value
// 	return nil
// }

var _commandParamReg = regexp.MustCompile(`\:\w+`)

func (cmd *DBCommand) Build() (sql string, params []interface{}, err error) {
	if !cmd.rebuild {
		sql = cmd.cachedSQL
		params = cmd.cachedParams
		err = nil
		return
	}
	sql = cmd.template
	params = make([]interface{}, 0)
	for {
		loc := _commandParamReg.FindStringIndex(sql)
		if len(loc) < 1 {
			break
		}

		name := strings.ToLower(sql[loc[0]+1 : loc[1]])
		if val, ok := cmd.parameters[name]; ok {
			params = append(params, val)
		} else {
			err = logs.NewError("DBERR_UNSETPARAM", "[DB] Named parameter :"+name+" not found, forgot set it?")
			return
		}
		sql = sql[:loc[0]] + "?" + sql[loc[1]:]
	}
	if sql == "" {
		err = logs.NewError("DBERR_BUILD", "SQL statement is required. template:"+ cmd.template)
			return
	}
	cmd.rebuild = false
	cmd.cachedSQL = sql
	cmd.cachedParams = params
	return
}

// type ExecutableCommand struct {
// 	*Command
// }

func (cmd *DBCommand) Exec(conn Connection) (sql.Result, error) {
	query, args, err := cmd.Build()
	if err != nil {
		return nil, logs.Wrap(err)
	}
	logs.Debug("\n[SQL] %v \n[Parameters] %+v", query, args)
	return conn.Exec(query, args...)
}

// func newExecutable() *ExecutableCommand {
// 	return &ExecutableCommand{
// 		Command: initCommand(),
// 	}
// }

// type QueryableCommand struct {
// 	*Command
// }

func (cmd *DBCommand) Query(conn Connection) (*sql.Rows, error) {
	query, args, err := cmd.Build()
	if err != nil {
		return nil, logs.Wrap(err)
	}
	logs.Debug("\n[SQL       ] %v \n[Parameters] %+v", query, args)
	return conn.Query(query, args...)
	// stmt, err := conn.Prepare(query)
	// if err != nil {
	// 	return nil, logs.Wrap(err)
	// }
	// defer stmt.Close()
	// return stmt.Query(args)
}

// func newQueryable(sql string) *QueryableCommand {
// 	return &QueryableCommand{
// 		Command: initCommand(sql),
// 	}
// }

// type DBCommand struct {
// 	*Command
// 	*ExecutableCommand
// 	*QueryableCommand
// }

// func NewCommand(sql string) *DBCommand {
// 	cmd := initCommand(sql)
// 	return &DBCommand{
// 		Command: cmd,
// 		QueryableCommand: &QueryableCommand{
// 			Command: cmd,
// 		},
// 		ExecutableCommand: &ExecutableCommand{
// 			Command: cmd,
// 		},
// 	}
// }
