// Package sqlite
//
// @author: xwc1125
package sqlite

type DBSql struct {
	Sql  string        `json:"sql" mapstructure:"sql"`
	Args []interface{} `json:"args" mapstructure:"args"`
}
