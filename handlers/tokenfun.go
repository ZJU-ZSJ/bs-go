package handlers

import (
	"bs-go/database"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"strconv"
	"time"
)

func tokencreate(username string) string {
	h := md5.New()
	_, _ = io.WriteString(h, username)
	_, _ = io.WriteString(h, time.Now().String())
	res := h.Sum(nil)
	log.Printf(hex.EncodeToString(res))
	return hex.EncodeToString(res)
}

func tokenrefresh(username string) string {
	tokenTmp := tokencreate(username)
	now := time.Now()
	m, _ := time.ParseDuration("10m")
	stmt, err := database.DBCon.Prepare("update User set token=?,time=? where username=?")
	if err != nil {
		log.Printf("%q", err)
	}
	_, err = stmt.Exec(tokenTmp, now.Add(m), username)
	if err != nil {
		log.Printf("%q", err)
	}
	return tokenTmp
}

func tokeninvalid(username string) bool {
	res := true
	var time_tmp time.Time
	err := database.DBCon.QueryRow("SELECT time FROM User WHERE username=?", username).Scan(&time_tmp)
	if err != nil {
		log.Printf("%q", err)
	}
	if time_tmp.Unix() < time.Now().Unix() {
		res = false
	}
	return res
}

func didlogin(uid int, token string) bool {
	if uid == 0 || len(token) == 0 {
		return false
	}
	res := true
	var tokenTmp string
	err := database.DBCon.QueryRow("SELECT token FROM User WHERE uid=?", strconv.Itoa(uid)).Scan(&tokenTmp)
	if err != nil {
		log.Printf("%q", err)
	}
	if token != tokenTmp {
		res = false
	}
	return res
}
