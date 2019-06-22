package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Handlewant(c *gin.Context) {
	i, err := c.Request.Cookie("uid")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "未登录！",
		})
		return
	}
	toke, err := c.Request.Cookie("token")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "未登录！",
		})
		return
	}
	uid, _ := strconv.Atoi(i.Value)
	token := toke.Value
	if !didlogin(uid, token) {
		c.JSON(http.StatusOK, gin.H{
			"code": -10,
			"msg":  "未登录",
		})
		return
	}

	wantedid, _ := strconv.Atoi(c.Request.PostFormValue("wantedid"))

	var bookname string
	var pricewanted float64
	var buyerid int

	row := database.DBCon.QueryRow("SELECT bookname,pricewanted,uid FROM BookWanted WHERE wantedid=?", wantedid)
	err = row.Scan(&bookname, &pricewanted, &buyerid)
	if buyerid == uid {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "你不能提供自己的求购！",
		})
		return
	}

	returncode := 0
	var id = int64(0)
	stmt, err := database.DBCon.Prepare("INSERT INTO Ord(bookid,bookname, ordertime,buyerid,salerid,ordertype,state,bcom,scom) values(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		res, err := stmt.Exec(0, bookname, time.Now(), buyerid, uid, 0, 0, 0, 0)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "订单提交失败！",
			})
			return
		}
		id, _ = res.LastInsertId()
	}
	if returncode == 0 {
		_, err := database.DBCon.Exec(
			"update BookWanted set state=?,orderid=? where wantedid=?",
			1,
			id,
			wantedid,
		)
		if err != nil {
			stmt, err = database.DBCon.Prepare("delete from Ord where ordid=?")
			_, err = stmt.Exec(id)
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "订单提交失败！",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"data": gin.H{
				"msg": "订单提交成功！",
				"id":  id,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "订单提交失败！",
		})
	}
}
