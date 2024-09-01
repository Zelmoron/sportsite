package main

import (
	"net/http"
	"site/handlers"
)

func main() {

	go handlers.Handlers()
	http.ListenAndServe(":8080", nil)

}
