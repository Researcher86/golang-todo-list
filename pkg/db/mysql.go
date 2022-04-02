package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	db *sql.DB
}

func NewMySQL(dsn string, maxOpenConns int, maxIdleConns int) *MySQL {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	return &MySQL{db: db}
}

func (mysql *MySQL) Close() {
	err := mysql.db.Close()
	if err != nil {
		panic(err)
	}
}

func (mysql *MySQL) Select(query string) (*sql.Rows, error) {
	results, err := mysql.db.Query(query)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (mysql *MySQL) Insert(query string) (int64, error) {
	res, err := mysql.db.Exec(query)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (mysql *MySQL) Update(query string) (int64, error) {
	res, err := mysql.db.Exec(query)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}
