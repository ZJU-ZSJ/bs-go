package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func LoginPage(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")
	log.Println(username, password)

	returncode := 0

	var password_tmp string
	var token_tmp string
	var time_tmp time.Time
	var uid int

	err := database.DBCon.QueryRow("SELECT uid,password,token,time FROM User WHERE username=?", username).Scan(&uid, &password_tmp, &token_tmp, &time_tmp)
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		if password != password_tmp {
			returncode = -2
		}
		log.Printf(time_tmp.String())
		if time_tmp.Unix() < time.Now().Unix() {
			token_tmp = tokenrefresh(username)
		}
		log.Printf(token_tmp)
	}
	if returncode == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"data": gin.H{
				"token":    token_tmp,
				"username": username,
				"uid":      uid,
			},
		})
	} else if returncode == -1 {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "该用户未注册",
		})
	} else if returncode == -2 {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "密码错误",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "未知错误",
		})
	}
}
