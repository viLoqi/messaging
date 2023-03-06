package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type WriteResponseBody struct {
	MessageId string `json:"messageID,omitempty"`
}

type ReadResponseBody struct {
	Author      string "author"
	Content     string "content"
	TimeCreated int    "timeCreated"
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// func TestReadFireStoreHandlerFound(t *testing.T) {
// 	requestBody, _ := json.Marshal(map[string]string{
// 		"fullMessagePath": "chats/SAM101/sec01/room/messages/CvJ0pNJw8HZ6izKARGfW",
// 	})

// 	responseBody, _ := json.Marshal(map[string]interface{}{
// 		"author":      "Testing Script",
// 		"content":     "This is to test read functionality",
// 		"timeCreated": 1677634337,
// 	})

// 	r := SetUpRouter()
// 	r.GET("/api/messaging", ReadFireStoreHandler)
// 	req, _ := http.NewRequest("GET", "/api/messaging", bytes.NewBuffer(requestBody))
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	responseData, _ := io.ReadAll(w.Body)
// 	assert.Equal(t, responseBody, responseData)
// 	assert.Equal(t, http.StatusOK, w.Code)
// }

// func TestReadFireStoreHandlerNotFound(t *testing.T) {
// 	requestBody, _ := json.Marshal(map[string]string{
// 		"fullMessagePath": "chats/SAM101/sec01/room/messages/aksdbabskdbasdbk",
// 	})

// 	responseBody, _ := json.Marshal(map[string]interface{}{
// 		"fullMessagePath": "chats/SAM101/sec01/room/messages/aksdbabskdbasdbk",
// 	})

// 	r := SetUpRouter()
// 	r.GET("/api/messaging", ReadFireStoreHandler)
// 	req, _ := http.NewRequest("GET", "/api/messaging", bytes.NewBuffer(requestBody))
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	responseData, _ := io.ReadAll(w.Body)
// 	assert.Equal(t, responseBody, responseData)
// 	assert.Equal(t, http.StatusNotFound, w.Code)
// }

func TestReadWriteAndDeleteFireStoreHandler(t *testing.T) {
	var __write WriteResponseBody
	var __read ReadResponseBody
	var __expectedRead ReadResponseBody

	postRequestBody, _ := json.Marshal(map[string]string{
		"collectionPath": "chats/SAM101/sec01/room/messages",
		"content":        "This is to test read functionality",
	})

	b := []byte(`{"author":"Testing Script","content":"This is to test read functionality", "timeCreated": 12312312}`)
	json.Unmarshal(b, &__expectedRead)

	r := SetUpRouter()
	r.GET("/api/messaging", ReadFireStoreHandler)
	r.POST("/api/messaging", WriteFireStoreHandler)
	r.DELETE("/api/messaging", DeleteFireStoreHandler)

	postRequest, _ := http.NewRequest("POST", "/api/messaging", bytes.NewBuffer(postRequestBody))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, postRequest)

	_rb, _ := io.ReadAll(w.Body)

	json.Unmarshal(_rb, &__write)

	getRequestBody, _ := json.Marshal(map[string]string{
		"fullMessagePath": "chats/SAM101/sec01/room/messages/" + __write.MessageId,
	})

	getRequest, _ := http.NewRequest("GET", "/api/messaging", bytes.NewBuffer(getRequestBody))

	r.ServeHTTP(w, getRequest)

	_grb, _ := io.ReadAll(w.Body)

	json.Unmarshal(_grb, &__read)

	assert.Equal(t, __read.Author, __expectedRead.Author)
	assert.Equal(t, http.StatusOK, w.Code)
}

// func TestDeleteFireStoreHandlerNotFound(t *testing.T) {
// 	testCollectionPath := "chats/SAM101/sec01/room/messages"
// 	messageID := "kdurYFGaEk5ggGau8v43"

// 	fullMessagePath := fmt.Sprintf("%s/%s", testCollectionPath, messageID)

// 	requestBody, _ := json.Marshal(map[string]string{
// 		"fullMessagePath": fullMessagePath,
// 	})

// 	responseBody, _ := json.Marshal(map[string]interface{}{
// 		"removed": fullMessagePath,
// 	})

// 	r := SetUpRouter()
// 	r.DELETE("/api/messaging", DeleteFireStoreHandler)
// 	req, _ := http.NewRequest("DELETE", "/api/messaging", bytes.NewBuffer(requestBody))
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	responseData, _ := io.ReadAll(w.Body)
// 	assert.Equal(t, responseBody, responseData)
// 	assert.Equal(t, http.StatusOK, w.Code)

// }
