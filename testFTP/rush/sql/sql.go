package main

import (
	gosql "database/sql"
	"fmt"
	"net/url"

	//Register sql drivers
	_ "github.com/denisenkom/go-mssqldb" // mssql (sql server)
	_ "github.com/go-sql-driver/mysql"   // mysql
	_ "github.com/jackc/pgx/v4/stdlib"   // pgx (postgres)
	_ "github.com/mattn/go-sqlite3"      // sqlite
	_ "github.com/sijms/go-ora/v2"       // oracle
)

type Plugin struct {
	Type     string            `yaml:"type"`
	Url      string            `yaml:"url"`
	Path     string            `yaml:"path"`
	Format   string            `yaml:"format"`
	User     string            `yaml:"user"`
	PassWord string            `yaml:"password"`
	Method   string            `yaml:"method"`
	Headers  map[string]string `yaml:"headers"`
}

var config = &Plugin{
	Type:     "postgres",
	Url:      "127.0.0.1:5432",
	Path:     "rush_ts030",
	Format:   "",
	User:     "admin",
	PassWord: "admin",
	Method:   "",
	Headers:  map[string]string{},
}

var items = map[string]string{"1": "hello ftp", "2": "kevin"}

const (
	PostgresProviderType  = "postgres"
	MysqlProviderType     = "mysql"
	SqlServerProviderType = "sqlserver"
	OracleProviderType    = "oracle"
	SqliteProviderType    = "sqlite"

	defaultTableName    = "rush_publish_package"
	TableExistsTemplate = "SELECT 1 FROM %s"
	CreateTableTemplate = `create table rush_publish_package
	(
		id              %s not null
				primary key,
		controller_name varchar(255),
		tool_sn         varchar(255),
		update_time     varchar(255),
		step_data       text,
		count           integer,
		pset            integer,
		job             integer,
		nut_no          varchar(255),
		batch           varchar(255),
		vin          varchar(255),
		tightening_id   varchar(255),
		curve_data      text
	);`
	InsertIntoTemplate = `insert into rush_publish_package (controller_name,
tool_sn,update_time,step_data,count,pset,job,nut_no,batch,vin,tightening_id,curve_data) 
values ('%s','%s','%s','%s',%d,%d,%d,'%s','%s','%s','%s','%s')`
)

var (
	driverMap = map[string]string{
		"postgres":  "pgx",
		"sqlserver": "mssql",
		"mysql":     "mysql",
		"sqlite":    "sqlite",
	}
)

type SQL struct {
	Driver         string
	Type           string
	dataSourceName string

	db *gosql.DB
}

func NewSQLProvider(config *Plugin) *SQL {
	var dataSourceName string
	driver, ok := driverMap[config.Type]
	if !ok {
		return nil
	}
	switch config.Type {
	case PostgresProviderType:
		dataSourceName = fmt.Sprintf("%s://%s:%s@%s/%s", config.Type, config.User, config.PassWord, config.Url, config.Path)
	case MysqlProviderType:
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s)/%s", config.User, config.PassWord, config.Url, config.Path)
	case SqlServerProviderType:
		query := url.Values{}
		query.Add("database", config.Path)
		u := &url.URL{
			Scheme:   "sqlserver",
			User:     url.UserPassword(config.User, config.PassWord),
			Host:     config.Url,
			RawQuery: query.Encode(),
		}
		dataSourceName = u.String()
	case SqliteProviderType:
		dataSourceName = config.Path
	}
	return &SQL{
		Type:           config.Type,
		Driver:         driver,
		dataSourceName: dataSourceName,
	}
}

func (s *SQL) generateInsert(pkg string) string {
	return fmt.Sprintf(
		InsertIntoTemplate,
		pkg,
		pkg,
		pkg,
		pkg,
		1,
		2,
		3,
		pkg,
		pkg,
		pkg,
		pkg,
		pkg,
	)
}

func (s *SQL) tableExists(tableName string) bool {
	sql := fmt.Sprintf(TableExistsTemplate, tableName)
	_, err := s.db.Exec(sql)
	return err == nil
}

func (s *SQL) generateCreateTable() string {
	var identity string
	switch s.Type {
	case SqlServerProviderType:
		identity = "integer identity(1,1)"
	case SqliteProviderType:
		identity = "integer"
	default:
		identity = "serial"
	}
	return fmt.Sprintf(CreateTableTemplate, identity)
}

func (s *SQL) Connect() error {
	db, err := gosql.Open(s.Driver, s.dataSourceName)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	s.db = db

	if !s.tableExists(defaultTableName) {
		_, err := s.db.Exec(s.generateCreateTable())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQL) Write(data map[string]string) error {
	for _, pkg := range data {
		sql := s.generateInsert(pkg)
		_, err := s.db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQL) Close() error {
	return s.db.Close()
}

func main() {
	f := NewSQLProvider(config)
	_ = f.Connect()

	_ = f.Write(items)
	_ = f.Close()
}
