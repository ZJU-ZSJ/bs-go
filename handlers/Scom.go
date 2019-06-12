package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func Scom(c *gin.Context) {
	orderid := c.PostForm("orderid")
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
	var bcom int
	var scom int
	err := database.DBCon.QueryRow("SELECT bcom,scom,salerid FROM Ord WHERE orderid=?", orderid).Scan(&bcom, &scom, &salerid)
	if err != nil {
		log.Printf("%q", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "订单号错误！",
		})
		return
	}
	if salerid != uid {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "无权限操作此订单！",
		})
		return
	}
	if bcom == 1 && scom == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"msg": "确认发货成功！",
			},
		})
	}
	state := 0
	if bcom == 1 {
		state = 1
	}
	stmt, err := database.DBCon.Prepare("update Ord set scom=?,state=? where orderid=?")

	_, err = stmt.Exec(1, state, orderid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "确认发货失败！",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"msg": "确认发货成功！",
		},
	})
}
