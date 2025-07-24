package main

import (
	"gotodo/routers"
	"net/http"
	"os"
)

// imports the multiplexer and puts it into the listen and serve
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r := routers.InitRoutes()
	http.ListenAndServe(":"+port, r) 
}
