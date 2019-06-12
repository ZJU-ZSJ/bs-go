package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func Userinfo(c *gin.Context) {
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

	var username string
	var email string
	var name string
	var face string
	returncode := 0
	err := database.DBCon.QueryRow("SELECT username,email,name,face FROM User WHERE uid=?", uid).Scan(&username, &email, &name, &face)
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"data": gin.H{
				"username": username,
				"email":    email,
				"name":     name,
				"face":     face,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "查询错误！",
		})
	}
}
