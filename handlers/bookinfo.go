package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Bookinfo(c *gin.Context) {
	reserve := c.DefaultQuery("reserve", "0")
	offset := c.DefaultQuery("offset", "0")
	max := c.DefaultQuery("num", "50")
	returncode := 0
	var querystr string
	if reserve == "1" {
		querystr = "SELECT bookid,bookname,pricenow,pic,content FROM Book order by -bookid LIMIT " + offset + "," + max
	} else {
		querystr = "SELECT bookid,bookname,pricenow,pic,content FROM Book LIMIT " + offset + "," + max
	}

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
	for rows.Next() {
		var bookid int
		var bookname string
		var pricenow float32
		var pic string
		var content string
		err = rows.Scan(&bookid, &bookname, &pricenow, &pic, &content)
		if err != nil {
			log.Printf("%q", err)
		}
		data = append(data, gin.H{
			"bookid":   bookid,
			"bookname": bookname,
			"pricenow": pricenow,
			"pic":      pic,
			"content":  content,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": returncode,
		"data": data,
	})
}
