package handlers

import (
	"bs-go/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Bookshow(c *gin.Context) {
	id := c.Param("id")
	var bookname string
	var priceori float32
	var pricenow float32
	var category string
	var content string
	var pic string
	var bookurl string
	returncode := 0
	err := database.DBCon.QueryRow("SELECT bookname,priceori,pricenow,category,content,pic,bookurl FROM Book WHERE bookid=?", id).Scan(&bookname, &priceori, &pricenow, &category, &content, &pic, &bookurl)
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
	}
	if returncode == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"data": gin.H{
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
