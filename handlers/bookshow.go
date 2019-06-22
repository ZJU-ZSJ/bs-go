package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Bookshow(c *gin.Context) {
	id := c.Param("id")
	var username string
	var name string
	var bookname string
	var priceori float32
	var pricenow float32
	var category string
	var content string
	var pic string
	var bookurl string
	var state int
	var uid int
	returncode := 0
	err := database.DBCon.QueryRow("SELECT username,name,Book.uid as uid,state,bookname,priceori,pricenow,category,content,pic,bookurl FROM Book inner join User on Book.uid=User.uid WHERE bookid=?", id).Scan(&username, &name, &uid, &state, &bookname, &priceori, &pricenow, &category, &content, &pic, &bookurl)
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"data": gin.H{
				"username": username,
				"name":     name,
				"uid":      uid,
				"state":    state,
				"bookname": bookname,
				"priceori": priceori,
				"pricenow": pricenow,
				"category": category,
				"content":  content,
				"pic":      pic,
				"bookurl":  bookurl,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "查无此书！",
		})
	}
}
