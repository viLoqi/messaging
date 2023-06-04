package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"loqi/messaging/structs"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

var unitTestCollection string = "chats/TEST/01-LEC/room/messages"
var testingRoute string = "/api/messaging"

func MakePostRequest(r *gin.Engine, w *httptest.ResponseRecorder, postRequestBody []byte) string {
	var writeResponseFromAPI structs.WriteResponseBody

	postRequest, _ := http.NewRequest("POST", testingRoute, bytes.NewBuffer(postRequestBody))
	r.ServeHTTP(w, postRequest)
	postResponseBody, _ := io.ReadAll(w.Body)
	json.Unmarshal(postResponseBody, &writeResponseFromAPI)

	return unitTestCollection + "/" + writeResponseFromAPI.MessageId
}

func MakeGetRequest(r *gin.Engine, w *httptest.ResponseRecorder, getRequestBody []byte) structs.ReadResponseBody {

	var readResponseFromAPI structs.ReadResponseBody

	getRequest, _ := http.NewRequest("GET", testingRoute, bytes.NewBuffer(getRequestBody))
	r.ServeHTTP(w, getRequest)
	getResponseBody, _ := io.ReadAll(w.Body)
	json.Unmarshal(getResponseBody, &readResponseFromAPI)

	return readResponseFromAPI
}

func MakeDeleteRequest(r *gin.Engine, w *httptest.ResponseRecorder, deleteRequestBody []byte) structs.DeleteResponseBody {

	var deleteResponseFromAPI structs.DeleteResponseBody

	deleteRequest, _ := http.NewRequest("DELETE", testingRoute, bytes.NewBuffer(deleteRequestBody))
	r.ServeHTTP(w, deleteRequest)
	deleteResponseBody, _ := io.ReadAll(w.Body)
	json.Unmarshal(deleteResponseBody, &deleteResponseFromAPI)

	return deleteResponseFromAPI
}
