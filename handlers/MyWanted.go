package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func MyWanted(c *gin.Context) {
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
	querystr := "SELECT wantedid,bookname,pricewanted,moreinfo,state,time,IFNULL(orderid,0) FROM BookWanted WHERE uid=" + id.Value
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
		var wantedid int
		var bookname string
		var pricewanted float64
		var moreinfo string
		var addtime time.Time
		var state int
		var orderid int
		err = rows.Scan(&wantedid, &bookname, &pricewanted, &moreinfo, &state, &addtime, &orderid)
		if err != nil {
			log.Printf("%q", err)
		}
		data = append(data, gin.H{
			"wantedid":    wantedid,
			"key":         index,
			"bookname":    bookname,
			"pricewanted": pricewanted,
			"moreinfo":    moreinfo,
			"state":       state,
			"time":        addtime,
			"orderid":     orderid,
		})
		index++
	}
	c.JSON(http.StatusOK, gin.H{
		"code": returncode,
		"data": data,
	})
}
