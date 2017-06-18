# build-RESTful-service-with-mongodb-in-golang
This is a demo, building a RESTful service with MongoDB in Golang.

## Install
glide install

## Run 
go run main.go

## TEST RESTful Service
### Add a new user
```shell
curl -X POST -H 'Content-Type: application/json' -d @body.json http://localhsot:1424/user

body.json
{
  "id": "5",
  "name": "李四",
  "age": 11
}
```

### Edit a user
```shell
curl -X PUT -H 'Content-Type: application/json' -d @body.json http://localhost:1424/user

body.json
{
  "id": "1",
  "title": "天一",
  "age": "-1"
}
```

### Query all users
```shell
curl http://localhost:1424/users
```

### Query a user
```shell
curl http://localhost:1424/user
```

### Delete a user
```shell
curl -X DELETE http://localhost:1424/user/1
```
