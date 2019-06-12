package handlers

import (
	"bs-go/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterPage(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	log.Println(username, password)

	if len(username) < 4 || len(password) < 6 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "格式错误",
		})
		return
	}

	returncode := 0

	stmt, err := database.DBCon.Prepare("INSERT INTO User(username, password) values(?,?)")
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		_, err = stmt.Exec(username, password)
		if err != nil {
			returncode = -1
			log.Printf("%q", err)
		}
	}
	if returncode == 0 {
		tokenTmp := tokenrefresh(username)
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"data": gin.H{
				"username": username,
				"token":    tokenTmp,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "注册失败",
		})
	}
}
