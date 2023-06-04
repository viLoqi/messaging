package main

import (
	"log"
	core "loqi/messaging/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReadFireStoreHandler(c *gin.Context) {
	type ReadMessageRequestBody struct {
		FullMessagePath string `json:"fullMessagePath"`
	}

	client, ctx := core.CreateFireStoreClient()
	defer client.Close()

	var requestBody ReadMessageRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Read RequestBody Error: %s\n", err)
	}

	data, err := core.GetSnapShotData(client, ctx, requestBody.FullMessagePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func WriteFireStoreHandler(c *gin.Context) {
	type WriteMessageRequestBody struct {
		CollectionPath string `json:"collectionPath"`
		Content        string `json:"content"`
	}

	client, ctx := core.CreateFireStoreClient()
	defer client.Close()

	var requestBody WriteMessageRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Write RequestBody Error: %s\n", err)
	}

	ref, err := core.CreateNewDocument(client, ctx, requestBody.CollectionPath, requestBody.Content)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"messageID": ref.ID})
	}
}

func PatchFireStoreHandler(c *gin.Context) {
	type PatchMessageRequestBody struct {
		FullMessagePath string `json:"fullMessagePath"`
		Content         string `json:"content"`
	}

	client, ctx := core.CreateFireStoreClient()
	defer client.Close()

	var requestBody PatchMessageRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Patch RequestBody Error: %s\n", err)
	}

	_, err := core.UpdateDocument(client, ctx, requestBody.FullMessagePath, requestBody.Content)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}

}

func DeleteFireStoreHandler(c *gin.Context) {
	type DeleteMessageRequestBody struct {
		FullMessagePath string `json:"fullMessagePath"`
	}

	client, ctx := core.CreateFireStoreClient()
	defer client.Close()

	var requestBody DeleteMessageRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Delete RequestBody Error: %s\n", err)
	}

	_, err := core.DeleteDocument(client, ctx, requestBody.FullMessagePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"removedFullMessagePath": requestBody.FullMessagePath})
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
