package main

import (
	"log"
	core "loqi/messaging/core"
	"loqi/messaging/structs"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func ReadFireStoreHandler(c *gin.Context) {

	client, ctx := core.CreateFireStoreClient()
	defer client.Close()

	var requestBody structs.ReadMessageRequestBody
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

	client, ctx := core.CreateFireStoreClient()
	defer client.Close()

	var requestBody structs.WriteMessageRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Write RequestBody Error: %s\n", err)
	}

	// Need to improve this function so that it's more flexible to new fields
	ref, err := core.CreateNewDocument(client, ctx, requestBody)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"messageID": ref.ID})
	}
}

func PatchFireStoreHandler(c *gin.Context) {

	client, ctx := core.CreateFireStoreClient()
	defer client.Close()

	var requestBody structs.PatchMessageRequestBody
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

	client, ctx := core.CreateFireStoreClient()
	defer client.Close()

	var requestBody structs.DeleteMessageRequestBody
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

func AboutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Loqi's Messaging API"})
}

func main() {
	r := gin.Default()

	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	var pathPrefix = "/api/messaging"

	r.GET(pathPrefix+"/about", AboutHandler)
	r.GET(pathPrefix, ReadFireStoreHandler)
	r.POST(pathPrefix, WriteFireStoreHandler)
	r.DELETE(pathPrefix, DeleteFireStoreHandler)
	r.PATCH(pathPrefix, PatchFireStoreHandler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
