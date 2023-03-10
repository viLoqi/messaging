package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type ReadMessageRequestBody struct {
	FullMessagePath string `json:"fullMessagePath"`
}

type WriteMessageRequestBody struct {
	CollectionPath string `json:"collectionPath"`
	Content        string `json:"content"`
}

type PatchMessageRequestBody struct {
	FullMessagePath string `json:"fullMessagePath"`
	Content         string `json:"content"`
}

func SanitizeFirestorePath(path string) string {
	if len(path) > 0 {
		lastIndex := len(path) - 1
		if path[lastIndex] == '/' {
			return path[:lastIndex]
		}
	}
	return path
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

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
		log.Printf("An error has occurred with the requestBody: %s", err)
	}

	path := SanitizeFirestorePath(requestBody.FullMessagePath)

	dsnap, err := client.Doc(path).Get(ctx)

	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	} else {
		m := dsnap.Data()
		c.JSON(http.StatusOK, m)
	}
}

func WriteFireStoreHandler(c *gin.Context) {
	var requestBody WriteMessageRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
		log.Printf("An error has occurred with the requestBody: %s", err)
	}

	client, ctx := CreateFireStoreClient()
	defer client.Close()

	collectionPath := SanitizeFirestorePath(requestBody.CollectionPath)

	ref := client.Collection(collectionPath).NewDoc()

	if _, err := ref.Set(ctx, gin.H{
		"author":       "Testing Script",
		"content":      requestBody.Content,
		"firstCreated": firestore.ServerTimestamp,
		"lastUpdated":  firestore.ServerTimestamp,
	}); err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred with writing the document: %s", err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"messageID": ref.ID,
		})
	}
}

func DeleteFireStoreHandler(c *gin.Context) {
	var requestBody ReadMessageRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
		log.Printf("An error has occurred with the requestBody: %s", err)
	}

	client, ctx := CreateFireStoreClient()
	defer client.Close()

	path := SanitizeFirestorePath(requestBody.FullMessagePath)

	if _, err := client.Doc(requestBody.FullMessagePath).Delete(ctx); err != nil {
		// Usually document not found
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"removedFullMessagePath": path,
		})
	}
}

func PatchFireStoreHandler(c *gin.Context) {
	var requestBody PatchMessageRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
		log.Printf("An error has occurred with the requestBody: %s", err)
	}

	client, ctx := CreateFireStoreClient()
	defer client.Close()

	path := SanitizeFirestorePath(requestBody.FullMessagePath)
	newContent := requestBody.Content

	if _, err := client.Doc(path).Update(ctx, []firestore.Update{
		{
			Path:  "content",
			Value: newContent,
		},
		{
			Path:  "lastUpdated",
			Value: firestore.ServerTimestamp,
		},
	}); err != nil {
		// Usually document not found
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"patchedFullMessagePath": path,
		})
	}

}

func main() {
	r := gin.Default()

	r.GET("/api/messaging", ReadFireStoreHandler)
	r.POST("/api/messaging", WriteFireStoreHandler)
	r.DELETE("/api/messaging", DeleteFireStoreHandler)
	r.PATCH("/api/messaging", PatchFireStoreHandler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
