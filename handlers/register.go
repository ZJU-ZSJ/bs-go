package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterPage(c *gin.Context) {
	username := c.Request.PostFormValue("username")
	password := c.Request.PostFormValue("password")
	email := c.Request.PostFormValue("email")

	log.Println(username, password, email)

	if len(username) < 6 || len(username) > 20 || len(password) < 6 || len(password) > 20 || len(email) < 3 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "格式错误",
		})
		return
	}

	returncode := 0

	stmt, err := database.DBCon.Prepare("INSERT INTO User(username, password,email,name,face) values(?,?,?,?,?)")
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		_, err = stmt.Exec(username, password, email, "", "")
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
