package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetMoveRequestFixture() MoveRequest {
	return MoveRequest{
		"test-id-1",
		[]Point{{0, 0}, {0, 1}},
		"test-game-id",
		10,
		10,
		0,
		[]Snake{
			Snake{
				[]Point{{0, 8}, {0, 9}},
				100,
				"test-id-1",
				"test-name-1",
				"",
			},
			Snake{
				[]Point{{2, 9}, {2, 8}},
				100,
				"test-id-2",
				"test-name-2",
				"",
			},
		},
	}
}

func TestFindClosestFood(t *testing.T) {
	//Arrange
	moveRequest := GetMoveRequestFixture()
	activeGame := GameFactory(
		&moveRequest,
	)

	//Act
	index := activeGame.findClosestFood()

	//Assert
	assert.Equal(t, Point{0, 1}, index[0].coord)

}

func TestSortedDirections(t *testing.T) {
	sd := sortedDirections{{"left", 1}, {"right", 0}, {"up", 1}, {"down", 1}}
	sort.Sort(sd)
	assert.Equal(t, 0, sd[0].heuristic)
	assert.Equal(t, "right", sd[0].direction)
}

func TestBackTrack_FindPath(t *testing.T) {
	//Arrange
	moveRequest := GetMoveRequestFixture()
	moveRequest.Snakes[0].Coords = []Point{{0, 9}, {0, 8}}
	activeGame := GameFactory(&moveRequest)
	bt := backTrack{
		[]string{"up", "down", "left", "right"},
		0,
	}

	//Act
	s := state{
		[]string{},
		"",
		activeGame.board,
		activeGame.head,
	}
	direction := bt.FindPath(s, activeGame.head, Point{0, 1})

	//Assert
	assert.Equal(t, []string{"left", "up", "up", "up", "up", "up", "up", "up", "up", "right"}, direction)
}

func TestBackTrack_FindPath2(t *testing.T) {
	//Arrange
	moveRequest := GetMoveRequestFixture()
	moveRequest.Snakes[0].Coords = []Point{{1, 0}, {0, 1}, {1, 1}}
	moveRequest.Snakes[1].Coords = []Point{{3, 1}, {3, 2}, {3, 3}}
	moveRequest.Food[0].X = 0
	moveRequest.Food[0].Y = 8
	activeGame := GameFactory(&moveRequest)
	s := state{
		[]string{},
		"",
		activeGame.board,
		Point{1, 0},
	}
	bt := backTrack{
		[]string{"up", "down", "left", "right"},
		0,
	}

	//Act
	path := bt.FindPath(s, s.head, moveRequest.Food[0])

	//Assert
	assert.Equal(t, []string{"left", "left", "down", "down", "down", "down", "down", "down", "down", "down", "right"}, path)
}

func TestState_HasEscapePath(t *testing.T) {
	//Arrange
	state := state{
		[]string{"up", "down"},
		"down",
		[][]int{
			{1, 1, 1, 1},
			{1, 0, 0, 1},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
		},
		Point{2, 1},
	}

	//Act
	ret := state.HasEscapePath("left")

	//Assert
	assert.Equal(t, false, ret)
}

func TestCloneState(t *testing.T) {
	//Arrange
	state := state{
		[]string{"up", "down"},
		"down",
		[][]int{
			{0, 0, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		Point{0, 2},
	}

	//Act
	newState := state.Clone()
	state.head.Y = 1
	state.board[0][1] = 1

	//Assert
	assert.Equal(t, 2, newState.head.Y)
	assert.Equal(t, 0, newState.board[0][1])
	assert.Equal(t, []string{}, newState.usedDirections)
	assert.Equal(t, "", newState.currentDirection)
}
