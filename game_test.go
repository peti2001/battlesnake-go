package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//type MoveRequest struct {
//	You    string  `json:"you"`
//	Food   []Point `json:"food"`
//	GameId string  `json:"game_id"`
//	Height int     `json:"height"`
//	Width  int     `json:"width"`
//	Turn   int     `json:"turn"`
//	Snakes []Snake `json:"snakes"`
//}

func TestFindClosestFood(t *testing.T) {
	//Arrange
	activeGame := GameFactory(
		&MoveRequest{
			"test-id-1",
			[]Point{{0,0}, {0,1}},
			"test-game-id",
			10,
			10,
			0,
			[]Snake{
				Snake{
					[]Point{{0,8}, {0,9}},
					100,
					"test-id-1",
					"test-name-1",
					"",
				},
				Snake{
					[]Point{{1, 9}, {1, 1}},
					100,
					"test-id-2",
					"test-name-2",
					"",
				},
			},
		},
	)

	//Act
	index := activeGame.findClosestFood()

	//Assert
	assert.Equal(t, 1, index)

}
