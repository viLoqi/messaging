# messaging

Write/Update records in Firestore

**Simulate API Requests with curl**

```bash
curl -X POST http://localhost:8080/api/messaging -H 'Content-Type: application/json' -d \
'{
    "collectionPath": "chats/SAM101/sec01/room/messages",
    "content": "This is written RIGHT NOW!"
}'

Output:
{
    "messageID": "ID of the new message"
}
```

```bash
curl -X GET http://localhost:8080/api/messaging -H 'Content-Type: application/json' -d \
'{
    "fullMessagePath": "chats/SAM101/sec01/room/messages/{messageID}"
}'

Output:
{
  "author": string,
  "content": string,
  "firstCreated": timestamp,
  "lastUpdated": timestamp
}
```

```bash
curl -X PATCH http://localhost:8080/api/messaging -H 'Content-Type: application/json' -d \
'{
    "fullMessagePath": "chats/SAM101/sec01/room/messages/{messageID}",
    "content": "This is a revised version"
}

Output:
{
    "patchedFullMessagePath": "Copy of fullMessagePath"
}
```

```bash
curl -X DELETE http://localhost:8080/api/messaging -H 'Content-Type: application/json' -d \
'{
    "fullMessagePath": "chats/SAM101/sec01/room/messages/{messageID}"
}'

Output:
{
    "removedFullMessagePath": "Copy of fullMessagePath"
}
```
