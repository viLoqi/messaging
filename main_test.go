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
	Author      string `json:"author,omitempty"`
	Content     string `json:"content,omitempty"`
	TimeCreated int    `json:"timeCreated,omitempty"`
}

type DeleteResponseBody struct {
	RemovedFullMessagePath string `json:"removedFullMessagePath,omitempty"`
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestReadWriteAndDeleteFireStoreHandler(t *testing.T) {
	//Set Up
	var writeResponseFromAPI WriteResponseBody
	var readResponseFromAPI ReadResponseBody
	var deleteResponseFromAPI DeleteResponseBody
	var expectedReadResponseFromAPI ReadResponseBody

	// Initialization
	unitTestCollection := "chats/SAM101/sec01/room/messages"
	json.Unmarshal([]byte(`{"author":"Testing Script","content":"This is to test read functionality", "timeCreated": 12312312}`), &expectedReadResponseFromAPI)

	// Setting up Routes
	r := SetUpRouter()
	r.GET("/api/messaging", ReadFireStoreHandler)
	r.POST("/api/messaging", WriteFireStoreHandler)
	r.DELETE("/api/messaging", DeleteFireStoreHandler)
	w := httptest.NewRecorder()

	// Testing POST Functionalities
	postRequestBody, _ := json.Marshal(map[string]string{
		"collectionPath": unitTestCollection,
		"content":        "This is to test read functionality",
	})
	postRequest, _ := http.NewRequest("POST", "/api/messaging", bytes.NewBuffer(postRequestBody))

	r.ServeHTTP(w, postRequest)
	postResponseBody, _ := io.ReadAll(w.Body)
	json.Unmarshal(postResponseBody, &writeResponseFromAPI)
	testMessageRefPath := unitTestCollection + "/" + writeResponseFromAPI.MessageId

	assert.Equal(t, http.StatusOK, w.Code)

	// Testing GET Functionalities
	getRequestBody, _ := json.Marshal(map[string]string{
		"fullMessagePath": testMessageRefPath,
	})
	getRequest, _ := http.NewRequest("GET", "/api/messaging", bytes.NewBuffer(getRequestBody))

	r.ServeHTTP(w, getRequest)
	getResponseBody, _ := io.ReadAll(w.Body)
	json.Unmarshal(getResponseBody, &readResponseFromAPI)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedReadResponseFromAPI.Author, readResponseFromAPI.Author)
	assert.Equal(t, expectedReadResponseFromAPI.Content, readResponseFromAPI.Content)

	// Testing DELETE Functionalities
	expectedDeleteResponseFromAPI := testMessageRefPath

	deleteRequestBody, _ := json.Marshal(map[string]string{
		"fullMessagePath": testMessageRefPath,
	})
	deleteRequest, _ := http.NewRequest("DELETE", "/api/messaging", bytes.NewBuffer(deleteRequestBody))

	r.ServeHTTP(w, deleteRequest)
	deleteResponseBody, _ := io.ReadAll(w.Body)
	json.Unmarshal(deleteResponseBody, &deleteResponseFromAPI)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedDeleteResponseFromAPI, deleteResponseFromAPI.RemovedFullMessagePath)
}

func TestReadFireStoreHandlerNotFound(t *testing.T) {
	requestBody, _ := json.Marshal(map[string]string{
		"fullMessagePath": "invalidPath",
	})

	expectedResponseBody := "\"firestore: nil DocumentRef\""

	r := SetUpRouter()
	r.GET("/api/messaging", ReadFireStoreHandler)
	req, _ := http.NewRequest("GET", "/api/messaging", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, expectedResponseBody, string(responseData))
	assert.Equal(t, http.StatusNotFound, w.Code)
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
