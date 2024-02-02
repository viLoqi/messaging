# messaging

- The current implementation requires the `credentials.json` file

Service to write/update/delete records stored inside Firebase Firestore.

**Simulate API Requests with curl**


**POST**
```bash
curl -X POST http://localhost:8080/api/messaging -H 'Content-Type: application/json' -d \
'{
    "collectionPath": "chats/SAM101/sec01/room/messages",
    "content": "This is written RIGHT NOW!",
    "author": "Jie Chen",
    "authorPhotoURL": "www."
}'
```
Response
```tsx
{
    "messageID": uuid
}
```

**GET**
```bash
curl -X GET http://localhost:8080/api/messaging -H 'Content-Type: application/json' -d \
'{
    "fullMessagePath": "chats/SAM101/sec01/room/messages/{messageID}"
}'
```
Response
```tsx
{
    "author": string,
    "content": string,
    "firstCreated": Date,
    "lastUpdated": Date
}
```

**PATCH**
```bash
curl -X PATCH http://localhost:8080/api/messaging -H 'Content-Type: application/json' -d \
'{
    "fullMessagePath": "chats/SAM101/sec01/room/messages/{messageID}",
    "content": "This is a revised version"
}
```
Response
```tsx
{
    "UpdateTime": Date
}
```

**DELETE**
```bash
curl -X DELETE http://localhost:8080/api/messaging -H 'Content-Type: application/json' -d \
'{
    "fullMessagePath": "chats/SAM101/sec01/room/messages/{messageID}"
}'
```
Response
```tsx
{
    "removedFullMessagePath": string
}
```
