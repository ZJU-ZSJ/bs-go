package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func Changename(c *gin.Context) {
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
	newname := c.Request.PostFormValue("newname")
	log.Printf(newname)
	if len(newname) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "昵称格式错误",
		})
		return
	}
	stmt, err := database.DBCon.Prepare("update User set name=? where uid=?")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更换昵称失败！",
		})
		return
	}
	_, err = stmt.Exec(newname, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更换昵称失败！",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"msg": "更换昵称成功！",
		},
	})
}
