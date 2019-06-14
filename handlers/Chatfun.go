package handlers

import (
	"bsgo/database"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func chatinit(uid int, aid int) int {
	fmt.Println("uid:", uid, "aid:", aid)
	var id int
	myid := 0
	anid := 0
	_ = database.DBCon.QueryRow("SELECT id FROM Chat WHERE user_id=? and another_id=?", uid, aid).Scan(&myid)
	_ = database.DBCon.QueryRow("SELECT id FROM Chat WHERE user_id=? and another_id=?", aid, uid).Scan(&anid)
	if myid != 0 {
		return myid
	} else if anid != 0 {
		return anid
	} else {
		stmt, err := database.DBCon.Prepare("INSERT INTO Chat(user_id,another_id) values(?,?)")
		if err != nil {
			fmt.Println(err)
			return -1
		}
		res, err := stmt.Exec(uid, aid)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		id_tmp, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
			return -1
		}
		id = int(id_tmp)
		stmt, err = database.DBCon.Prepare("INSERT INTO ChatList(chat_id,user_id,another_id,is_online) values(?,?,?,?)")
		if err != nil {
			fmt.Println(err)
			return -1
		}
		res, err = stmt.Exec(id, uid, aid, 1)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		res, err = stmt.Exec(id, aid, uid, 0)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		fmt.Println("chatid:", id)
		stmt, err = database.DBCon.Prepare("INSERT INTO ChatMsg(chat_id,user_id,content,time) values(?,?,?,?);\n")
		if err != nil {
			fmt.Println(err)
			return -1
		}
		_, err = stmt.Exec(id, uid, "hello", time.Now())
		if err != nil {
			fmt.Println(err)
			return -1
		}
		return id
	}
}

func chatinto(id int, uid int) {
	fmt.Println("chatinto", id, "uid:", uid)
	stmt, err := database.DBCon.Prepare("update ChatList set is_online=1,unread=0 where chat_id=? and user_id=?")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(id, uid)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func chatget(chat_id int, offset int, limit int) gin.H {
	querystr := "SELECT content,time,user_id FROM ChatMsg WHERE chat_id=" + strconv.Itoa(chat_id) + " order by id desc LIMIT " + strconv.Itoa(offset) + "," + strconv.Itoa(limit)
	fmt.Println(querystr)
	rows, err := database.DBCon.Query(querystr)
	if err != nil {
		fmt.Println(err)
		return gin.H{"code": -1}
	}
	var data = []gin.H{}
	for rows.Next() {
		var content string
		var time time.Time
		var user_id int
		err = rows.Scan(&content, &time, &user_id)
		if err != nil {
			log.Printf("%q", err)
		}
		data = append(data, gin.H{
			"content": content,
			"time":    time,
			"user_id": user_id,
		})
	}
	return gin.H{
		"code": 0,
		"data": data,
	}
}

func chatend(id int, uid int) {
	fmt.Println("chatout", id)
	stmt, err := database.DBCon.Prepare("update ChatList set is_online=0 where chat_id=? and user_id=?")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(id, uid)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("byebye")
	return
}

func chatsend(chatlistid int, uid int, content string) gin.H {
	fmt.Println("uid:", uid, "chatlistid:", chatlistid)
	now := time.Now()
	var aonline int
	_ = database.DBCon.QueryRow("SELECT is_online FROM ChatList WHERE another_id=? and chat_id=?", uid, chatlistid).Scan(&aonline)
	if aonline == 0 {
		stmt, err := database.DBCon.Prepare("update ChatList set unread=unread+1 where another_id=? and chat_id=?")
		if err != nil {
			fmt.Println(err)
			return gin.H{"code": -1}
		}
		_, err = stmt.Exec(uid, chatlistid)
		if err != nil {
			fmt.Println(err)
			return gin.H{"code": -1}
		}
	}
	stmt, err := database.DBCon.Prepare("UPDATE ChatMsg set is_latest=0 WHERE id IN (SELECT id FROM ChatMsg WHERE is_latest=1 and user_id=? and chat_id=? ORDER BY id DESC LIMIT 1)")
	if err != nil {
		fmt.Println(err)
		return gin.H{"code": -1}
	}
	_, err = stmt.Exec(uid, chatlistid)
	if err != nil {
		fmt.Println(err)
		return gin.H{"code": -1}
	}
	stmt, err = database.DBCon.Prepare("INSERT INTO ChatMsg(chat_id,user_id,content,time) values(?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return gin.H{"code": -1}
	}
	_, err = stmt.Exec(chatlistid, uid, content, now)
	if err != nil {
		fmt.Println(err)
		return gin.H{"code": -1}
	}
	var data = []gin.H{}
	data = append(data, gin.H{
		"content": content,
		"time":    now,
		"user_id": uid,
	})
	return gin.H{
		"code": 0,
		"data": data,
	}
}
