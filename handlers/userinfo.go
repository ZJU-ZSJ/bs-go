package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Userinfo(c *gin.Context) {
	uid := c.Param("uid")

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
