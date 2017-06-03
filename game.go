package main

import (
	"math"
	"fmt"
)

type Game struct {
	me             string
	width, height  int
	foods []Point
	head Point
}

func GameFactory(me string, width, height int, foods []Point, head Point) *Game {
	return &Game{
		me,
		width,
		height,
		foods,
		head,
	}
}

func (g *Game) findClosestFood() int {
	minIndex := 0
	minDistance := math.Pow(float64(g.width + g.height), 2)

	for i, food := range g.foods {
		if d := math.Pow(float64(food.X - g.head.X), 2) + math.Pow(float64(food.Y - g.head.Y), 2); d < minDistance {
			minDistance = d
			minIndex = i
		}
	}

	return minIndex
}

func (g *Game) ChooseDirection() string {
	i := g.findClosestFood()
	fmt.Println("Closest: ", i)
	if dX := g.foods[i].X - g.head.X; dX > 0 {
		return "right"
	} else if dX < 0 {
		return "left"
	}
	if dY := g.foods[i].Y - g.head.Y; dY > 0 {
		return "down"
	} else if dY < 0 {
		return "up"
	}
	return "up"
}