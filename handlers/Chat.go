package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[int]*websocket.Conn)

type data struct {
	Action  string `json:"action"`
	Uid     string `json:"uid"`
	Token   string `json:"token"`
	Msg     string `json:"msg"`
	Aneroid string `json:"aid"`
}

func Chat(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	login := false
	chatlistid := -1
	uid := -1
	defer ws.Close()
	for {
		// 读取ws中的数据
		_, message, err := ws.ReadMessage()
		if err != nil {
			// 客户端关闭连接时也会进入
			fmt.Println(err)
			if login {
				delete(clients, uid)
				chatend(chatlistid, uid)
			}
			break
		}
		msg := &data{}
		_ = json.Unmarshal(message, msg)
		fmt.Println(string(message))

		if msg.Action == "join" {
			uid, err = strconv.Atoi(msg.Uid)
			if err != nil {
				break
			}
			if !didlogin(uid, msg.Token) {
				v := gin.H{
					"code": -10,
					"msg":  "未登录",
				}
				err = ws.WriteJSON(v)
				if err != nil {
					break
				}
				break
			}
			login = true
			fmt.Println(msg.Uid, "joined")
			clients[uid] = ws
			uid = uid
		} else if msg.Action == "start" {
			aid, _ := strconv.Atoi(msg.Aneroid)
			if err != nil {
				break
			}
			chatlistid = chatinit(uid, aid)
			if err != nil {
				break
			}
			chatinto(chatlistid, uid)
			clients[uid].WriteJSON(chatget(chatlistid, 0, 100))
		} else if msg.Action == "send" {
			aid, _ := strconv.Atoi(msg.Aneroid)
			v := chatsend(chatlistid, uid, msg.Msg)
			if _, ok := clients[aid]; ok {
				_ = clients[aid].WriteJSON(v)
			}
			v = gin.H{"message": "已发送"}
			_ = ws.WriteJSON(v)

		}
	}
}
