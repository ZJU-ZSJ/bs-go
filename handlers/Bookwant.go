package handlers

import (
	"bsgo/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func BookWant(c *gin.Context) {
	i, _ := c.Request.Cookie("uid")
	toke, _ := c.Request.Cookie("token")
	uid, _ := strconv.Atoi(i.Value)
	token := toke.Value
	if !didlogin(uid, token) {
		c.JSON(http.StatusOK, gin.H{
			"code": -10,
			"msg":  "未登录",
		})
		return
	}
	bookname := c.Request.PostFormValue("bookname")
	pricewanted, _ := strconv.ParseFloat(c.Request.PostFormValue("pricewanted"), 64)
	moreinfo := c.Request.PostFormValue("moreinfo")
	log.Printf("%s,%f,%s", bookname, pricewanted, moreinfo)

	if len(bookname) == 0 || len(moreinfo) >= 500 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "格式错误",
		})
		return
	}
	returncode := 0
	var id = int64(0)

	stmt, err := database.DBCon.Prepare("INSERT INTO BookWanted(bookname, pricewanted,moreinfo,uid,state,time) values(?,?,?,?,?,?)")
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		res, err := stmt.Exec(bookname, pricewanted, moreinfo, uid, 0, time.Now())
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"code": returncode,
				"msg":  "求购添加失败",
			})
			return
		}
		id, _ = res.LastInsertId()
	}
	if returncode == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"data": gin.H{
				"msg": "求购添加成功",
				"id":  id,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "求购添加失败",
		})
	}
}
