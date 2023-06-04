package main

import (
	"encoding/json"
	"loqi/messaging/structs"
	h "loqi/messaging/testing"
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

// Initialization
var unitTestCollection string = "chats/TEST/01-LEC/room/messages"
var testingRoute string = "/api/messaging"

func TestReadWriteAndDeleteFireStoreHandler(t *testing.T) {
	//Set Up
	var expectedReadResponseFromAPI structs.ReadResponseBody

	// Initialization
	json.Unmarshal([]byte(`{"author":"Testing Script","content":"This is to test read functionality", "firstCreated": 12312312}`), &expectedReadResponseFromAPI)

	// Setting up Routes
	r := SetUpRouter()
	r.GET(testingRoute, ReadFireStoreHandler)
	r.POST(testingRoute, WriteFireStoreHandler)
	r.DELETE(testingRoute, DeleteFireStoreHandler)
	w := httptest.NewRecorder()

	// Testing POST Functionalities
	postRequestBody, _ := json.Marshal(map[string]string{
		"collectionPath": unitTestCollection,
		"content":        "This is to test read functionality",
	})

	testMessageRefPath := h.MakePostRequest(r, w, postRequestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	// Testing GET Functionalities
	getRequestBody, _ := json.Marshal(map[string]string{
		"fullMessagePath": testMessageRefPath,
	})

	readResponseFromAPI := h.MakeGetRequest(r, w, getRequestBody)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedReadResponseFromAPI.Author, readResponseFromAPI.Author)
	assert.Equal(t, expectedReadResponseFromAPI.Content, readResponseFromAPI.Content)

	// Testing DELETE Functionalities
	expectedDeleteResponseFromAPI := testMessageRefPath

	deleteRequestBody, _ := json.Marshal(map[string]string{
		"fullMessagePath": testMessageRefPath,
	})

	deleteResponseFromAPI := h.MakeDeleteRequest(r, w, deleteRequestBody)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedDeleteResponseFromAPI, deleteResponseFromAPI.RemovedFullMessagePath)
}
