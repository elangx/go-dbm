package dbm

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var connections map[string]*sql.DB

func init() {
	connections = make(map[string]*sql.DB)
}

type DBConfig struct {
	host     string
	user     string
	password string
	dbname   string
	port     int
	setting  []Setting
}

const (
	defaultPort    = 3306
	defaultTimeout = "1000"
	sDSNFormat     = "%s%s=%s"
)

type Setting func(string) string

func stringSetting(source string, key string, value string) string {
	if value == "" {
		return source
	}
	return fmt.Sprintf(sDSNFormat, source, key, value)
}

func boolSetting(source string, key string, ok bool) string {
	return fmt.Sprintf(sDSNFormat, source, key, strconv.FormatBool(ok))
}

func timeSetting(source string, key string, t time.Duration) string {
	if t < time.Millisecond || t >= 24*time.Hour {
		return source
	}
	return fmt.Sprintf(sDSNFormat, source, key, t)
}

func SetTimeout(t time.Duration) Setting {
	return func(source string) string {
		return timeSetting(source, "timeout", t)
	}
}

func SetCharset(v string) Setting {
	return func(source string) string {
		return stringSetting(source, "charset", v)
	}
}

func New(host string, user string, password string, dbname string) *DBConfig {
	return &DBConfig{
		host:     host,
		user:     user,
		password: password,
		dbname:   dbname,
		port:     defaultPort,
	}
}

func (p *DBConfig) Port(port int) *DBConfig {
	p.port = port
	return p
}

func (p *DBConfig) Set(sets ...Setting) *DBConfig {
	p.setting = append(p.setting, sets...)
	return p
}

func (p *DBConfig) Add(tagName string) error {
	db_ins, err := sql.Open("mysql", p.getRealDSN())
	if err != nil {
		return err
	}
	err = db_ins.Ping()
	if err != nil {
		return err
	}
	connections[tagName] = db_ins
	return nil
}

func (p *DBConfig) getRealDSN() string {
	dsn := "%s:%s@tcp(%s:%d)/%s?%s"
	return fmt.Sprintf(dsn, p.user, p.password, p.host, p.port, p.dbname, concatDSN(p.setting))
}

func concatDSN(sets []Setting) string {
	s := ""
	for _, f := range sets {
		s = f(s)
	}
	return s
}

func GetConnection(tagName string) (*sql.DB, error) {
	if _, ok := connections[tagName]; ok {
		return connections[tagName], nil
	}
	return nil, fmt.Errorf("connection %s: does not exists", tagName)
}
