package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func BookAdd(c *gin.Context) {
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
	priceori := c.Request.PostFormValue("priceori")
	pricenow := c.Request.PostFormValue("pricenow")
	category := c.Request.PostFormValue("category")
	content := c.Request.PostFormValue("content")
	pic := c.Request.PostFormValue("pic")
	bookurl := c.Request.PostFormValue("bookurl")
	log.Printf("%s,%f,%f,%s,%s,%s,%s",bookname,priceori,pricenow,category,content,pic,bookurl)

	if len(bookname) == 0 || len(priceori) == 0 || len(pricenow) == 0 || len(category) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "格式错误",
		})
		return
	}
	returncode := 0
	var id = int64(0)

	stmt, err := database.DBCon.Prepare("INSERT INTO Book(bookname, priceori,pricenow,category,content,pic,bookurl,uid,state,time) values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		res, err := stmt.Exec(bookname, priceori, pricenow, category, content, pic, bookurl, uid, 0, time.Now())
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": returncode,
				"msg":  "图书添加失败",
			})
			return
		}
		id, _ = res.LastInsertId()
	}
	if returncode == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"data": gin.H{
				"msg": "图书添加成功",
				"id":  id,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "图书添加失败",
		})
	}
}
