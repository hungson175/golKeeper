//Thiis package is used for models & data store
package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
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

func NewDataSource() (*DataSource, error) {
	mysqlUsername := os.Getenv("MYSQL_USERNAME")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/golkeeper", mysqlUsername, mysqlPassword))
	return &DataSource{db}, err
}
