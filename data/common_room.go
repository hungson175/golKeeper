//Package data is used for models & data store
package data

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql" //just sql, dont call init
)

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

func newDatabase() *sql.DB {
	mysqlUsername := os.Getenv(EnviVarSQLUsername)
	mysqlPassword := os.Getenv(EnvVarSQLPassword)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUsername, mysqlPassword, IPAddress, Port, dbName))
	if err != nil {
		panic(err)
	}
	return db
}

var dbSync sync.Once
var database *sql.DB

func GetDBInstance() *sql.DB {
	dbSync.Do(func() {
		fmt.Println("Database Init: begin....")
		database = newDatabase()
		fmt.Println("Database init: done ")
	})
	if database == nil {
		panic(errors.New("nil db returned"))
	}
	return database
}
