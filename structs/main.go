package structs

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

type DeleteMessageRequestBody struct {
	FullMessagePath string `json:"fullMessagePath"`
}

type ReadResponseBody struct {
	Author       string `json:"author"`
	Content      string `json:"content"`
	FirstCreated int    `json:"firstCreated"`
	LastUpdated  int    `json:"lastUpdated"`
}

type WriteResponseBody struct {
	MessageId string `json:"messageID"`
}

type DeleteResponseBody struct {
	RemovedFullMessagePath string `json:"removedFullMessagePath"`
}
