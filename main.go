package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kananb/blind-chess/controllers/api"
	"github.com/kananb/blind-chess/controllers/game"
)

func main() {
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// route API
	group := router.Group("/api")
	api.Route(group)

	// route game
	group = router.Group("/game")
	game.Route(group)

	router.Run()
}
