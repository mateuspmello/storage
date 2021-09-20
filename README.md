# API Upload file

API Upload file is a service that provide you upload files and download thoses files 

## What does API Upload Files

API Upload Files is the most easy way to upload file to a server.


## Prerequisites

You will need the following tools properly installed on your computer.

* [Git](http://git-scm.com/)
* [Go](http://golang.org/)

## Running Storage test

```shell
go test storagedata/storagedata_test.go
```
## Running API test

```shell
go test storagedata/storagedata_test.go
```

## Running the service

```shell
go run cmd/apiamericanas/main.go
```

## End Points

## Send file
    
POST /sendfile
### Curl example:
```bash
curl \
-F "file=@/home/mateus-mello/go/src/americanas/test_files/curiosity.png " \
https://localhost:8081/sendfile \
--data '{"path": "fromcurl"}' -v
```

## Delete file
    
POST /delete?data=FileID
### Curl example:
```bash
curl -X POST 'http://localhost:8081/delete?data=fa0ecd5f42635c34e2f879a24039988e'
```

## All files
    
GET /allfiles
### Curl example:
```bash
curl -X POST 'http://localhost:8081/allfiles'
```
    
    GET /storagedata/*filepath

## Get file by id
    
GET /byid?data=FileID
### Curl example:
```bash
curl -X POST 'http://localhost:8081/byid?data=fa0ecd5f42635c34e2f879a24039988e'
```

## Move file to another dir
    
 POST /movefile?data=FileID
### Curl example:
```bash
curl -X POST 'http://localhost:8081/movefile?data=fa0ecd5f42635c34e2f879a24039988e'
```

## Overwrite file
    
POST /overwrite?data=FileID
### Curl example:
```bash
curl -X POST 'http://localhost:8081/overwrite?data=fa0ecd5f42635c34e2f879a24039988e'
```

## Download file
    
GET /storagedata/dirOfFile
### Curl example:
```bash
curl -X GET 'http://localhost:8081/delete?data=fa0ecd5f42635c34e2f879a24039988e'
```

    