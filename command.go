package db

import (
	"database/sql"
	"regexp"
	"strings"

	"github.com/one-go/logs"
)

type Command struct {
	conn         Connection
	parameters   map[string]interface{}
	template     string
	rebuild      bool // built
	cachedSQL    string
	cachedParams []interface{}
}

func NewCommand(conn Connection, template ...string) *Command {
	cmd := &Command{
		conn:       conn,
		parameters: make(map[string]interface{}),
		rebuild:    true,
	}
	if len(template) > 0 {
		cmd.SetCommand(template[0])
	}
	return cmd
}
func (cmd *Command) SetCommand(template string) error {
	cmd.template = template
	cmd.rebuild = true
	return nil
}

func (cmd *Command) SetParam(name string, value interface{}) error {
	cmd.parameters[strings.ToLower(name)] = value
	return nil
}

var _commandParamReg = regexp.MustCompile(`\:\w+`)

func (cmd *Command) Build() (sql string, params []interface{}, err error) {
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
	cmd.rebuild = false
	cmd.cachedSQL = sql
	cmd.cachedParams = params
	return
}

func (cmd *Command) Exec() (sql.Result, error) {
	query, args, err := cmd.Build()
	if err != nil {
		return nil, err
	}
	logs.Debug("\n[SQL] %v \n[PARAMS]%+v", query, args)
	return cmd.conn.Exec(query, args...)
}

func (cmd *Command) Query() (*sql.Rows, error) {
	query, args, err := cmd.Build()
	if err != nil {
		return nil, err
	}
	logs.Debug("\n[SQL] %v \n[Parameters]%+v", query, args)
	return cmd.conn.Query(query, args...)
}
