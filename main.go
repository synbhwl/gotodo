package main

import (
	"gotodo/routers"
	"net/http"
)

// imports the multiplexer and puts it into the listen and serve 
func main() {
	r := routers.InitRoutes()
	http.ListenAndServe(":3000", r) 
}
