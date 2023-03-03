package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
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

func CreateFireStoreClient() (*firestore.Client, context.Context) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("./credentials.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client, ctx
}

func ReadFireStoreHandler(c *gin.Context) {
	var requestBody ReadMessageRequestBody

	client, ctx := CreateFireStoreClient()
	defer client.Close()

	if flag.Lookup("test.v") == nil {
		fmt.Println("normal run")
	} else {
		fmt.Println("run under go test")
	}

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
	}

	dsnap, err := client.Collection("chats").Doc("SAM101").Collection("sec01").Doc("room").Collection("messages").Doc(requestBody.MessageID).Get(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	m := dsnap.Data()
	fmt.Printf("Document data: %#v\n", m)

	c.JSON(http.StatusOK, m)
}

func WriteFireStoreHandler(c *gin.Context) {
	var requestBody WriteMessageRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
	}

	client, ctx := CreateFireStoreClient()
	defer client.Close()

	if flag.Lookup("test.v") == nil {
		fmt.Println("normal run")
	} else {
		fmt.Println("run under go test")
	}

	ref := client.Collection("chats").Doc("SAM101").Collection("sec01").Doc("room").Collection("messages").NewDoc()

	_, err := ref.Set(ctx, map[string]interface{}{
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

func DeleteFireStoreHandler(c *gin.Context) {
	var requestBody ReadMessageRequestBody

	if flag.Lookup("test.v") == nil {
		fmt.Println("normal run")
	} else {
		fmt.Println("run under go test")
	}

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
	}

	client, ctx := CreateFireStoreClient()
	defer client.Close()

	_, err := client.Collection("chats").Doc("SAM101").Collection("sec01").Doc("room").Collection("messages").Doc(requestBody.MessageID).Delete(ctx)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"removed": requestBody.MessageID,
	})

}

func main() {
	r := gin.Default()

	r.GET("/api/messaging", ReadFireStoreHandler)
	r.POST("/api/messaging", WriteFireStoreHandler)
	r.DELETE("/api/messaging", DeleteFireStoreHandler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
