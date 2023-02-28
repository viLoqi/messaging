package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type MessageRequestBody struct {
	MessageID string `json:"messageID,omitempty"`
}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Loqi's Messaging API"})
}

func ReadFireStoreHandler(c *gin.Context) {
	var requestBody MessageRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
	}

	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("./cred.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	dsnap, err := client.Collection("chats").Doc("SAM101").Collection("sec01").Doc("room").Collection("messages").Doc(requestBody.MessageID).Get(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	m := dsnap.Data()
	fmt.Printf("Document data: %#v\n", m)

	c.JSON(http.StatusOK, gin.H{
		"message": requestBody.MessageID,
	})
}

func main() {
	r := gin.Default()

	r.GET("/", HomepageHandler)
	r.GET("/read", ReadFireStoreHandler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
