package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type ReadMessageRequestBody struct {
	MessageID string `json:"messageID,omitempty"`
}

type WriteMessageRequestBody struct {
	Content string `json:"content,omitempty"`
}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Loqi's Messaging API"})
}

func ReadFireStoreHandler(c *gin.Context) {
	var requestBody ReadMessageRequestBody

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

func WriteFireStoreHandler(c *gin.Context) {
	var requestBody WriteMessageRequestBody

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

	ref := client.Collection("chats").Doc("SAM101").Collection("sec01").Doc("room").Collection("messages").NewDoc()

	_, err = ref.Set(ctx, map[string]interface{}{
		"author":      "Testing Script",
		"content":     requestBody.Content,
		"timeCreated": time.Now().Unix(),
	})
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"messageID": ref.ID,
	})

}

func main() {
	r := gin.Default()

	r.GET("/", HomepageHandler)
	r.GET("/read", ReadFireStoreHandler)
	r.POST("/write", WriteFireStoreHandler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
