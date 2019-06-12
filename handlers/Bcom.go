package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	_ "time"
)

func Bcom(c *gin.Context) {
	orderid := c.PostForm("orderid")
	log.Printf(orderid)
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
	var buyerid int
	var bcom int
	var scom int
	err := database.DBCon.QueryRow("SELECT bcom,scom,buyerid FROM Ord WHERE orderid=?", orderid).Scan(&bcom, &scom, &buyerid)
	if err != nil {
		log.Printf("%q", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "订单号错误！",
		})
		return
	}
	if buyerid != uid {
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
				"msg": "确认收货成功！",
			},
		})
	}
	var state = 0
	if scom == 1 {
		state = 1
	}
	stmt, err := database.DBCon.Prepare("update Ord set bcom=?,state=? where orderid=?")
	if err != nil {
		log.Printf("%q", err)
	}
	_, err = stmt.Exec(1, state, orderid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "确认收货失败！",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"msg": "确认收货成功！",
		},
	})
}
