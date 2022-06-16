package websocket

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gomokuAI/pkg/common"
	"log"
	"strconv"
	"strings"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type Message struct {
	SenderId string `json:"sender_id"`
	Type     int    `json:"type"`
	Body     string `json:"body"`
}

// 存储玩家与棋盘
var Player map[string][][]int

func (c *Client) Read() {
	defer func() {
		// 两个*才能更改c的pool的具体内容
		delete(Player, c.ID)
		fmt.Printf("当前玩家人数：%d\n", len(Player))
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		msg := string(p)
		tp, _ := strconv.Atoi(strings.Split(msg, " ")[0])
		tmp := strings.Split(msg, " ")[1:]
		msg = strings.Join(strings.Split(msg, " ")[1:], " ")
		if tp == 2 {
			fmt.Printf("%s 进来了\n", c.ID)
			chess := Player[c.ID]
			x, _ := strconv.Atoi(tmp[0])
			y, _ := strconv.Atoi(tmp[1])
			fmt.Printf("玩家的位置： %d %d\n", x, y)
			chess[x][y] = 1
			px, py := common.AI(chess)
			chess[px][py] = 2
			Player[c.ID] = chess
			fmt.Printf("AI的位置： %d %d\n", px, py)
			c.Conn.WriteJSON("2 " + strconv.Itoa(px) + " " + strconv.Itoa(py))
			fmt.Printf("%s 出去了\n", c.ID)
		}
		if tp == 3 {
			id := uuid.New().String()
			c.ID = id
			Player[c.ID] = common.InitChess()
			chess := Player[c.ID]
			fmt.Printf("%s玩家进来了\n", c.ID)
			fmt.Printf("当前玩家人数：%d\n", len(Player))
			if msg == "0" {
				c.Conn.WriteJSON(struct {
					body string
				}{
					body: "2 7 7",
				})
				chess[7][7] = 2
				Player[c.ID] = chess
			}
		}
		if tp == 66 {
			log.Printf("%s: %s\n", c.ID, msg)
		}
	}
}
