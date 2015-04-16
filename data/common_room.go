//Package data is used for models & data store
package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" //just sql, dont call init
)

type DataSource struct {
	db *sql.DB
}

type NotImplementedError struct {
	name string
}

func (err NotImplementedError) Error() string {
	return err.name + ": Not implemeted yet"
}

const IPAddress = "127.0.0.1"
const Port = "3306"
const EnviVarSQLUsername = "MYSQL_USERNAME"
const EnvVarSQLPassword = "MYSQL_PASSWORD"
const dbName = "golkeeper"

func NewDataSource() (*DataSource, error) {
	mysqlUsername := os.Getenv(EnviVarSQLUsername)
	mysqlPassword := os.Getenv(EnvVarSQLPassword)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUsername, mysqlPassword, IPAddress, Port, dbName))
	return &DataSource{db}, err
}
