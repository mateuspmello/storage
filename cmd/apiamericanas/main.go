package main

import (
	"americanas/api"
	"americanas/storagedata"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {

	storage := storagedata.New()
	router := httprouter.New()
	api.New(storage).RegisterRouters(router)
	fmt.Println("api Server running on http://localhost:8081")
	panic(http.ListenAndServe(":8081", router))
}
