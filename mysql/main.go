package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "xx:xx@tcp(localhost:3306)/eggsample?charset=utf8")
	defer db.Close()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(5 * time.Minute)
	checkErr(err)

	// insert
	stmt1, err := db.Prepare("INSERT user SET name=?,email=?")
	checkErr(err)
	defer stmt1.Close()
	res, err := stmt1.Exec("mysql", "mysql@bbb.com")
	checkErr(err)
	id, _ := res.LastInsertId()
	fmt.Println(id)
	// fmt.Println(res.RowsAffected())

	// update
	stmt2, err := db.Prepare("update user set name=? where id=?")
	checkErr(err)
	defer stmt2.Close()
	res, err = stmt2.Exec("coder", id)
	checkErr(err)

	// list
	rows, err := db.Query("SELECT id, name, email FROM user where id=?", id)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var email interface{}
		err = rows.Scan(&id, &name, &email)
		checkErr(err)
		fmt.Println(id, name, email)
	}

	// delete
	stmt3, err := db.Prepare("delete from user where id=?")
	checkErr(err)
	defer stmt3.Close()

	res, err = stmt3.Exec(id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		// panic(err)
	}
}
