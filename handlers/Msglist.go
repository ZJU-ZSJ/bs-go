package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func Msglist(c *gin.Context) {
	id, _ := c.Request.Cookie("uid")
	toke, _ := c.Request.Cookie("token")
	uid, _ := strconv.Atoi(id.Value)
	token := toke.Value
	if !didlogin(uid, token) {
		c.JSON(http.StatusOK, gin.H{
			"code": -10,
			"msg":  "未登录",
		})
		return
	}
	returncode := 0
	querystr := "SELECT another_id,unread FROM ChatList WHERE user_id=" + id.Value + " order by unread desc"
	rows, err := database.DBCon.Query(querystr)
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "查询出错！",
		})
		return
	}
	var data = []gin.H{}
	var index = 0
	for rows.Next() {
		var another_id int
		var unread int
		err = rows.Scan(&another_id, &unread)
		if err != nil {
			log.Printf("%q", err)
		}
		data = append(data, gin.H{
			"another_id": another_id,
			"unread":     unread,
		})
		index++
	}
	c.JSON(http.StatusOK, gin.H{
		"code": returncode,
		"data": data,
	})
}
