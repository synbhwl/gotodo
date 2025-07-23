package main

import (
	"gotodo/routers"
	"net/http"

)

func main(){
	r := routers.InitRoutes()
	http.ListenAndServe(":3000", r)
}