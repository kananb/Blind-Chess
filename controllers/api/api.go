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
		board, err = chess.BoardFromString(fen)
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

	// var strMoves []string
	// for move := range board.GenMoves() {
	// 	strMoves = append(strMoves, move.ToSAN(board))
	// }
	moves := board.GenMoves()
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

	move, err := chess.MoveFromString(c.Param("move"), board)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else if err := board.MakeMove(move); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"fen": board.FEN(),
		})
	}
}

func Route(router *gin.RouterGroup) {
	router.GET("/moves", getMoves)
	router.POST("/moves/:move", postMove)
}
