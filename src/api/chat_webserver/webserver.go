package main

import (
	"io"
	"log"
	"net/http"

	client "github.com/S-ign/chat-chat-chat/src/api/chat_client"
	"github.com/S-ign/chat-chat-chat/src/api/chat_webserver/room"
	"github.com/gin-gonic/gin"
)

var chatClient, chatConn = client.Connect()
var chatManager *room.Manager

func main() {
	chatManager = room.NewRoomManager()

	files := []string{
		"templates/index.html", "templates/chat.html",
	}

	r := gin.Default()

	r.Static("/js", "static/js")
	r.LoadHTMLFiles(files...) // load the built dist path
	r.GET("/chat/:roomid", chatGet)
	r.POST("/chat/:roomid", chatPost)
	r.DELETE("/chat/:roomid", chatDelete)
	r.GET("/stream/:roomid", stream)

	// Listen and serve on 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}
	defer chatConn.Close()
}

func stream(c *gin.Context) {
	roomid := c.Param("roomid")
	listener := chatManager.OpenListener(roomid)
	defer chatManager.CloseListener(roomid, listener)

	clientGone := c.Writer.CloseNotify()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			c.SSEvent("message", message)
			return true
		}
	})
}

func chatGet(c *gin.Context) {
	roomid := c.Param("roomid")
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"roomid": roomid,
	})
}

func chatPost(c *gin.Context) {
	roomid := c.Param("roomid")
	message := c.PostForm("chatmessage")
	if message != "" {
		chatManager.Submit(roomid, message)
		client.Chat(message, chatClient)

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": message,
		})
	}
}

func chatDelete(c *gin.Context) {
	roomid := c.Param("roomid")
	chatManager.DeleteBroadcast(roomid)
}
