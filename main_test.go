package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestReadFireStoreHandler(t *testing.T) {
	requestBody, _ := json.Marshal(map[string]string{
		"fullMessagePath": "chats/SAM101/sec01/room/messages/CvJ0pNJw8HZ6izKARGfW",
	})

	responseBody, _ := json.Marshal(map[string]interface{}{
		"author":      "Testing Script",
		"content":     "This is to test read functionality",
		"timeCreated": 1677634337,
	})

	r := SetUpRouter()
	r.GET("/api/messaging", ReadFireStoreHandler)
	req, _ := http.NewRequest("GET", "/api/messaging", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, responseBody, responseData)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWriteFireStoreHandler(t *testing.T) {
	requestBody, _ := json.Marshal(map[string]string{
		"fullMessagePath": "chats/SAM101/sec01/room/messages",
		"content":         "Wassup!!!",
	})

	r := SetUpRouter()
	r.POST("/api/messaging", WriteFireStoreHandler)
	req, _ := http.NewRequest("POST", "/api/messaging", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// responseData, _ := ioutil.ReadAll(w.Body)
	// assert.Equal(t, responseBody, responseData)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteFireStoreHandler(t *testing.T) {
	requestBody, _ := json.Marshal(map[string]string{
		"fullMessagePath": "chats/SAM101/sec01/room/messages",
	})

	r := SetUpRouter()
	r.DELETE("/api/messaging", DeleteFireStoreHandler)
	req, _ := http.NewRequest("DELETE", "/api/messaging", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// responseData, _ := ioutil.ReadAll(w.Body)
	// assert.Equal(t, responseBody, responseData)
	assert.Equal(t, http.StatusOK, w.Code)
}
