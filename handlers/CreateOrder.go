package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateOrder(c *gin.Context) {
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
	var salerid int
	bookid, _ := strconv.Atoi(c.Request.PostFormValue("bookid"))
	ordertype, _ := strconv.Atoi(c.Request.PostFormValue("ordertype"))
	address := c.Request.PostFormValue("address")

	if bookid <= 0 || ordertype < 0 || ordertype > 1 {
		log.Printf("%d,%d,%s", bookid, ordertype, address)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "格式错误",
		})
		return
	} else if ordertype == 0 && len(address) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "请填写地址",
		})
		return
	}
	var state int
	var bookname string
	err := database.DBCon.QueryRow("SELECT bookname,state,uid FROM Book WHERE bookid=?", bookid).Scan(&bookname, &state, &salerid)
	if err != nil {
		log.Printf("%q", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查无此书！",
		})
		return
	} else if uid == salerid {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "你不能购买自己的书",
		})
		return
	} else if state != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该书已被买走！",
		})
		return
	}

	returncode := 0
	var id = int64(0)
	stmt, err := database.DBCon.Prepare("INSERT INTO Ord(bookid,bookname, ordertime,buyerid,salerid,ordertype,address,state,bcom,scom) values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		res, err := stmt.Exec(bookid, bookname, time.Now(), uid, salerid, ordertype, address, 0, 0, 0)
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
		stmt, err = database.DBCon.Prepare("update Book set state=? where bookid=?")

		_, err = stmt.Exec(1, bookid)
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
