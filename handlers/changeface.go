package handlers

import (
	"bsgo/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Changeface(c *gin.Context) {
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
	imageUrl := c.PostForm("imageUrl")
	if len(imageUrl) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "图片链接错误",
		})
		return
	}
	stmt, err := database.DBCon.Prepare("update User set face=? where uid=?")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更换头像失败！",
		})
		return
	}
	_, err = stmt.Exec(imageUrl, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更换头像失败！",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"msg": "更换头像成功！",
		},
	})
}
