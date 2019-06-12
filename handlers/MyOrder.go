package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func MyOrder(c *gin.Context) {
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
	querystr := "SELECT bookid,orderid,bookname,ordertime,salerid,state,bcom,scom FROM Ord WHERE buyerid=" + id.Value
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
		var bookid int
		var orderid int
		var bookname string
		var ordertime time.Time
		var salerid int
		var state int
		var bcom int
		var scom int
		err = rows.Scan(&bookid, &orderid, &bookname, &ordertime, &salerid, &state, &bcom, &scom)
		if err != nil {
			log.Printf("%q", err)
		}
		data = append(data, gin.H{
			"bookid":    bookid,
			"key":       index,
			"orderid":   orderid,
			"bookname":  bookname,
			"ordertime": ordertime,
			"salerid":   salerid,
			"state":     state,
			"bcom":      bcom,
			"scom":      scom,
		})
		index++
	}
	c.JSON(http.StatusOK, gin.H{
		"code": returncode,
		"data": data,
	})
}
