package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func Wantedlist(c *gin.Context) {
	reserve := c.DefaultQuery("reserve", "0")
	offset := c.DefaultQuery("offset", "0")
	max := c.DefaultQuery("num", "50")
	returncode := 0
	var querystr string
	if reserve == "1" {
		querystr = "SELECT wantedid,bookname,pricewanted,moreinfo,state,time FROM BookWanted order by -wantedid LIMIT " + offset + "," + max
	} else {
		querystr = "SELECT wantedid,bookname,pricewanted,moreinfo,state,time FROM BookWanted LIMIT " + offset + "," + max

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
	var index = 0
	for rows.Next() {
		var wantedid int
		var bookname string
		var pricewanted float64
		var moreinfo string
		var addtime time.Time
		var state int
		err = rows.Scan(&wantedid, &bookname, &pricewanted, &moreinfo, &state, &addtime)
		if err != nil {
			log.Printf("%q", err)
		}
		data = append(data, gin.H{
			"wantedid":    wantedid,
			"key":         index,
			"bookname":    bookname,
			"pricewanted": pricewanted,
			"moreinfo":    moreinfo,
			"state":       state,
			"time":        addtime,
		})
		index++
	}
	c.JSON(http.StatusOK, gin.H{
		"code": returncode,
		"data": data,
	})
}
