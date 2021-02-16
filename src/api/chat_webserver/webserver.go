package main

import (
	client "github.com/S-ign/chat-chat-chat/src/api/chat_client"
	"github.com/gin-gonic/gin"
)

var chatClient = client.Connect()

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("chat_webserver/*.html")        // load the built dist path
	r.LoadHTMLFiles("static/*/*")                  //  load the static path
	r.Static("/static", "./chat_webserver/static") // use the loaded source
	r.StaticFile("/hello/", "dist/index.html")     // use the loaded source
	r.GET("/chat", chat)
	r.POST("/chat", chat)

	// Listen and serve on 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}
}

func chat(c *gin.Context) {
	chatClient.Chat(message, chatClient)
}
