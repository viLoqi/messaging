package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type MessageRequestBody struct {
	MessageID string `json:"messageID,omitempty"`
}

func main() {
	r := gin.Default()

	r.GET("/read", func(c *gin.Context) {

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

		iter := client.Collection("chats").Doc("testing").Collection("sec01").Documents(ctx)

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				break
			}
			fmt.Println(doc.Data())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": requestBody.MessageID,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
