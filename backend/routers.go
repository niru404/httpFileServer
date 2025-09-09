package routers

import (
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/gin-gonic/gin"
	"github.com/niru404/httpFileServer/backend/handlers"
)

func Setup() {
	router := gin.Default()

	router.GET("/ping", handlers.PingHandler)

	return router
}
