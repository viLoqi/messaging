package main

import (
	"log"
	h "loqi/messaging/helper"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func ReadFireStoreHandler(c *gin.Context) {
	type ReadMessageRequestBody struct {
		FullMessagePath string `json:"fullMessagePath"`
	}

	var requestBody ReadMessageRequestBody

	client, ctx := h.CreateFireStoreClient()
	defer client.Close()

	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Read RequestBody Error: %s\n", err)
	}

	path := requestBody.FullMessagePath

	data, err := h.GetSnapShotData(client, ctx, path)

	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func WriteFireStoreHandler(c *gin.Context) {
	type WriteMessageRequestBody struct {
		CollectionPath string `json:"collectionPath"`
		Content        string `json:"content"`
	}

	var requestBody WriteMessageRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Write RequestBody Error: \n%s", err)
	}

	client, ctx := h.CreateFireStoreClient()
	defer client.Close()

	path := requestBody.CollectionPath

	res, err := h.CreateNewDocument(client, ctx, path, requestBody.Content)

	if err != nil {
		log.Printf("Write Failure: %s\n", err)
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func DeleteFireStoreHandler(c *gin.Context) {
	type DeleteMessageRequestBody struct {
		FullMessagePath string `json:"fullMessagePath"`
	}
	var requestBody DeleteMessageRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
		log.Printf("Delete RequestBody Error: \n%s", err)
	}

	client, ctx := h.CreateFireStoreClient()
	defer client.Close()

	path := requestBody.FullMessagePath

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
	type PatchMessageRequestBody struct {
		FullMessagePath string `json:"fullMessagePath"`
		Content         string `json:"content"`
	}

	var requestBody PatchMessageRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
		log.Printf("Update RequestBody Error: \n%s", err)
	}

	client, ctx := h.CreateFireStoreClient()
	defer client.Close()

	path := requestBody.FullMessagePath
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

func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Loqi's Messaging API"})
}

func main() {
	r := gin.Default()

	var messagingServiceRoute = "/api/messaging"

	r.GET("/", HomeHandler)
	r.GET(messagingServiceRoute, ReadFireStoreHandler)
	r.POST(messagingServiceRoute, WriteFireStoreHandler)
	r.DELETE(messagingServiceRoute, DeleteFireStoreHandler)
	r.PATCH(messagingServiceRoute, PatchFireStoreHandler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
