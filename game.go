package main

import (
	"math"
)

type Game struct {
	data *MoveRequest
	head Point
	board         [][]int
}

func GameFactory(data *MoveRequest) *Game {
	var head Point
	board := make([][]int, data.Width)
	for i := 0;i < data.Width;i++ {
		board[i] = make([]int, data.Height)
	}

	for _, snake := range data.Snakes {
		if snake.Id == data.You {
			head.X = snake.Coords[0].X
			head.Y = snake.Coords[0].Y
		}
		for _, coord := range snake.Coords {
			board[coord.X][coord.Y] = 1
		}
	}

	return &Game{
		data,
		head,
		board,
	}
}

func (g *Game) findClosestFood() int {
	minIndex := 0
	minDistance := math.Pow(float64(g.data.Width + g.data.Height), 2)

	for i, food := range g.data.Food {
		if d := math.Pow(float64(food.X - g.head.X), 2) + math.Pow(float64(food.Y-g.head.Y), 2); d < minDistance {
			minDistance = d
			minIndex = i
		}
	}

	return minIndex
}

func (g *Game) ChooseDirection() string {
	i := g.findClosestFood()
	if dX := g.data.Food[i].X - g.head.X; dX > 0 {
		return "right"
	} else if dX < 0 {
		return "left"
	}
	if dY := g.data.Food[i].Y - g.head.Y; dY > 0 {
		return "down"
	} else if dY < 0 {
		return "up"
	}
	return "up"
}
