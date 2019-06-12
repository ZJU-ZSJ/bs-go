package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "strconv"
)

func Search(c *gin.Context) {
	//0:name
	//1:id
	//2:category
	searchtype := c.DefaultQuery("type", "0")
	content := c.DefaultQuery("content", "")
	returncode := 0
	querystr := ""
	if searchtype == "0" {
		querystr = "SELECT bookid,pic,bookname,category,state FROM Book WHERE bookname='" + content + "'"
	} else if searchtype == "1" {
		querystr = "SELECT bookid,pic,bookname,category,state FROM Book WHERE bookid=" + content
	} else if searchtype == "2" {
		querystr = "SELECT bookid,pic,bookname,category,state FROM Book WHERE category='" + content + "'"
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询出错！",
		})
		return
	}
	log.Printf(querystr)
	rows, err := database.DBCon.Query(querystr)
	if err != nil {
		returncode = -1
		log.Printf("%q", err)
		c.JSON(http.StatusOK, gin.H{
			"code": returncode,
			"msg":  "查询出错！",
		})
		return
	}
	var data = []gin.H{}
	var index = 0
	for rows.Next() {
		var bookid int
		var pic string
		var bookname string
		var category string
		var state int
		err = rows.Scan(&bookid, &pic, &bookname, &category, &state)
		if err != nil {
			log.Printf("%q", err)
		}
		data = append(data, gin.H{
			"bookid":   bookid,
			"pic":      pic,
			"bookname": bookname,
			"category": category,
			"state":    state,
			"key":      index,
		})
		index++
	}
	c.JSON(http.StatusOK, gin.H{
		"code": returncode,
		"data": data,
	})
}
