package main

import (
	"github.com/niru404/httpFileServer/backend"
)

func main() {
	router := routers.Setup()

	router.Run()
}


