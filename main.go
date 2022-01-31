package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kananb/blind-chess/controllers/api"
)

func main() {
	router := gin.Default()
	staticPath := "./frontend/blind-chess/build/"

	router.Use(static.Serve("/", static.LocalFile(staticPath, true)))

	// route API
	group := router.Group("/api")
	api.Route(group)

	router.LoadHTMLFiles(staticPath + "index.html")
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run()
}
