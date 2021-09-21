# API Upload file

API Upload file is a service that provide you upload files and download thoses files 

## What does API Upload Files

API Upload Files is the most easy way to upload file to a server.


## Prerequisites

You will need the following tools properly installed on your computer.

* [Git](http://git-scm.com/)
* [Go](http://golang.org/)

## Test

### Storage test

```shell
go test storagedata/storagedata_test.go
```
### API test

```shell
go test storagedata/storagedata_test.go
```

## Starting the service

```shell
go run cmd/apiamericanas/main.go
```

## End Points

### Send file
    
POST /sendfile
#### Curl example:
```bash
curl \
 -v -F path="ht/monthly" \
 -F file=@"/home/mateus-mello/go/src/americanas/test_files/mars.png" \
 http://localhost:8081/sendfile 

```

### All files
    
GET /allfiles
#### Curl example:
```bash
curl -X GET 'http://localhost:8081/allfiles'
```

### Delete file
    
POST /delete?data=FileID
#### Curl example:
```bash
curl -X POST 'http://localhost:8081/delete?data=fa0ecd5f42635c34e2f879a24039988e'
```

### Get file by id
    
GET /byid?data=FileID
#### Curl example:
```bash
curl -X GET 'http://localhost:8081/byid?data=0cb90ac871279cc942de976882b71a00'
```

### Move file to another dir
    
 POST /movefile?data=FileID
#### Curl example:
```bash
curl -X POST 'http://localhost:8081/movefile?data=0cb90ac871279cc942de976882b71a00' \
-d '{"directory": "solarsystem/planets"}'
```

### Overwrite file
    
POST /overwrite?data=FileID
#### Curl example:
```bash
curl \
 -F file=@"/home/mateus-mello/go/src/americanas/test_files/perseverance.png" \
 http://localhost:8081/overwrite?data=0cb90ac871279cc942de976882b71a00

```

### Download file
    
GET /storagedata/dirOfFile
#### Curl example:
```bash
curl -X GET 'http://localhost:8081/storagedata/solarsystem/planets/perseverance.png'
```
