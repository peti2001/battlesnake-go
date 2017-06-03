package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindClosestFood(t *testing.T) {
	//Arrange
	activeGame := GameFactory(
		"test",
		10,
		10,
		[]Point{{0,0}, {0,1}},
		Point{0, 10},
	)

	//Act
	index := activeGame.findClosestFood()

	//Assert
	assert.Equal(t, 1, index)

}
