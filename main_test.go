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
		"messageID": "CvJ0pNJw8HZ6izKARGfW",
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

// func TestWriteFireStoreHandler(t *testing.T) {
// 	mockResponse := `{"message":"Welcome to Loqi's Messaging API"}`
// 	r := SetUpRouter()
// 	r.POST("/api/messaging", HomepageHandler)
// 	req, _ := http.NewRequest("POST", "/api/messaging", nil)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	responseData, _ := ioutil.ReadAll(w.Body)
// 	assert.Equal(t, mockResponse, string(responseData))
// 	assert.Equal(t, http.StatusOK, w.Code)
// }
// func TestDeleteFireStoreHandler(t *testing.T) {
// 	mockResponse := `{"message":"Welcome to Loqi's Messaging API"}`
// 	r := SetUpRouter()
// 	r.DELETE("/api/messaging", HomepageHandler)
// 	req, _ := http.NewRequest("DELETE", "/api/messaging", nil)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	responseData, _ := ioutil.ReadAll(w.Body)
// 	assert.Equal(t, mockResponse, string(responseData))
// 	assert.Equal(t, http.StatusOK, w.Code)
// }
