package main

import (
	"bsgo/database"
	"bsgo/router"
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
			0：正常售卖
			1：已被买走

		订单状态:
			0: 已下单
			1：已发货或选择线下交易
			2：已完成

		订单分类：
			0:邮寄
			1:线下交易
	*/
	//创建表
	sql_table := `
    CREATE TABLE IF NOT EXISTS User(
        uid INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(64) NOT NULL UNIQUE,
        password VARCHAR(64) NOT NULL,
        email	 VARCHAR(64) NOT NULL UNIQUE,
        name 	 VARCHAR(64) NULL,
        token	 VARCHAR(64) NULL,
        face	VARCHAR(64) NULL,
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
    CREATE TABLE IF NOT EXISTS Ord(
        orderid INTEGER PRIMARY KEY AUTOINCREMENT,
        bookid  INTEGER NOT NULL,
        bookname VARCHAR NOT NULL,
        ordertime TIMESTAMP NOT NULL,
        buyerid INTEGER	NOT NULL,
        salerid INTEGER NOT NULL,
        ordertype INTEGER NOT NULL,
        address VARCHAR(64)  NULL,
        state	INTEGER NOT NULL,
        bcom	INTEGER NOT NULL,
        scom	INTEGER NOT NULL
    );
    `
	_, err = database.DBCon.Exec(sql_table)
	if err != nil {
		log.Printf("%q: %s\n", err, sql_table)
		return
	}
	router.Init() // init router
}
