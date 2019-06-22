package main

import (
	"bsgo/database"
	"bsgo/router"
	"database/sql"
	_ "database/sql"
	"flag"
	_ "fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	_ "time"
)

func main() {
	var err error
	database.DBCon, err = sql.Open("sqlite3", "./sql.db")
	if err != nil {
		log.Fatal(err)
	}
	baseurl := flag.String("url", "http://localhost:8080", "前端的地址")
	flag.Parse()
	log.Printf(*baseurl)
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

		求购分类:
			0:正常求购
			1:求购完成
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
	CREATE TABLE IF NOT EXISTS BookWanted(
        wantedid INTEGER PRIMARY KEY AUTOINCREMENT,
        bookname VARCHAR(64) NOT NULL,
        pricewanted FLOAT 		NOT NULL,
        moreinfo	 VARCHAR(64) NULL,
        uid		INTEGER NOT NULL,
        state	INTEGER NOT NULL,
        time	TIMESTAMP 	NOT NULL,
        orderid INTEGER NULL
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
    CREATE TABLE IF NOT EXISTS Chat(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        another_id INTEGER NOT NULL
    );
    CREATE TABLE IF NOT EXISTS ChatList(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        chat_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        another_id INTEGER NOT NULL,
        is_online INTEGER NOT NULL,
        unread INTEGER DEFAULT 0
    );
    CREATE TABLE IF NOT EXISTS ChatMsg(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        chat_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        content VARCHAR(64) NOT NULL,
        time 	TIMESTAMP NOT NULL
    );
    `
	_, err = database.DBCon.Exec(sql_table)
	if err != nil {
		log.Printf("%q: %s\n", err, sql_table)
		return
	}
	router.Init(*baseurl) // init router
}
