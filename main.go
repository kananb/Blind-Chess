package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kananb/blind-chess/controllers/game"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	staticPath := "./frontend/blind-chess/build/"

	router.Use(static.Serve("/", static.LocalFile(staticPath, true)))

	// route websockets
	group := router.Group("/game")
	game.Route(group)

	router.LoadHTMLFiles(staticPath + "index.html")
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run()
}
