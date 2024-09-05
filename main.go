package main

import (
	"fmt"
	"github.com/thinkeridea/go-extend/exnet"
	"gomokuAI/pkg/common"
	"gomokuAI/pkg/util"
	"gomokuAI/pkg/websocket"
	"log"
	"net/http"
	"strings"
)

func serveAI(w http.ResponseWriter, r *http.Request) {
	ip := exnet.ClientPublicIP(r)
	if ip == "" {
		ip = exnet.ClientIP(r)
	}
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	util.AskIpAddr(ip)
	log.Printf("wensocket enter")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
	}
	client.Read()
}

func setupRoutes() {
	common.InitAI()
	common.InitCache()
	websocket.Player = make(map[string][][]int)
	websocket.FirstPlayer = make(map[string]string)

	http.HandleFunc("/AI", func(w http.ResponseWriter, r *http.Request) {
		serveAI(w, r)
	})
	http.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request) {
		//ip := r.RemoteAddr
		ip := exnet.ClientPublicIP(r)
		if ip == "" {
			ip = exnet.ClientIP(r)
		}
		if strings.Contains(ip, ":") {
			ip = strings.Split(ip, ":")[0]
		}
		util.AskIpAddr(ip)
		w.Write([]byte("hello world\n"))
	})
}

func main() {
	fmt.Println("gomokuAI processing...")
	setupRoutes()
	_ = http.ListenAndServe(":9000", nil)
}
