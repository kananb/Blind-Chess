package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kananb/chess"
)

func initBoard(c *gin.Context) (board *chess.Board) {
	fen := c.Query("fen")

	var err error
	if fen != "" {
		board, err = chess.NewBoard(fen)
		if err != nil {
			c.Status(http.StatusBadRequest)
		}
	} else {
		board = chess.StartingPosition()
	}

	return board
}

func getMoves(c *gin.Context) {
	board := initBoard(c)
	if board == nil {
		return
	}

	var moves []string
	for _, move := range board.Moves() {
		moves = append(moves, move.String())
	}
	c.JSON(http.StatusOK, gin.H{
		"moves": moves,
		"count": len(moves),
	})
}

func postMove(c *gin.Context) {
	board := initBoard(c)
	if board == nil {
		return
	}

	move, err := chess.NewMove(c.Param("move"), board)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else if move = board.MakeMove(move); !move.IsValid() {
		c.String(http.StatusBadRequest, "Move is invalid")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"fen": board.String(),
			"san": move.String(),
		})
	}
}

func Route(router *gin.RouterGroup) {
	router.GET("/moves", getMoves)
	router.POST("/moves/:move", postMove)
}
