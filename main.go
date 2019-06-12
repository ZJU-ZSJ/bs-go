package main

import (
	"bs-go/database"
	"bs-go/router"
	"database/sql"
	_ "database/sql"
	_ "fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	_ "time"
)

func main() {
	var err error
	database.DBCon, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	/*
		书状态：
			1：正常售卖
			2：已被买走
	*/
	//创建表
	sql_table := `
    CREATE TABLE IF NOT EXISTS User(
        uid INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(64) NOT NULL UNIQUE,
        password VARCHAR(64) NOT NULL,
        email	 VARCHAR(64) NULL,
        name 	 VARCHAR(64) NULL,
        token	 VARCHAR(64) NULL,
        time 	 TIMESTAMP   NULL
    );
    CREATE TABLE IF NOT EXISTS Book(
        bookid INTEGER PRIMARY KEY AUTOINCREMENT,
        bookname VARCHAR(64) NOT NULL,
        priceori FLOAT 		NOT NULL,
        pricenow FLOAT 		NOT NULL,
        category VARCHAR(64) NOT NULL,
        content	 VARCHAR(64) NULL,
        pic 	 VARCHAR   NULL,
        bookurl	VARCHAR		NULL,
        uid		INTEGER NOT NULL,
        state	INTEGER NOT NULL,
        time	TIME 	NOT NULL
    );
    `
	_, err = database.DBCon.Exec(sql_table)
	if err != nil {
		log.Printf("%q: %s\n", err, sql_table)
		return
	}
	router.Init() // init router
}
