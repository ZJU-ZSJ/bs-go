package handlers

import (
	"bsgo/database"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Msgcount(c *gin.Context) {
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
	var unread int
	err := database.DBCon.QueryRow("SELECT SUM(unread) FROM ChatList WHERE user_id=?", uid).Scan(&unread)
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"unread": unread,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"unread": 0,
			},
		})
	}
}
