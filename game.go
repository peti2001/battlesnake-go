package main

import (
	"math"
	"sort"
	"fmt"
	"time"
)

type Game struct {
	data  *MoveRequest
	head  Point
	board [][]int
}

type state struct {
	usedDirections   []string
	currentDirection string
	board            [][]int
	head             Point
}

type backTrack struct {
	directions []string
}

type sortedDirections []directionWithHeuristic

type directionWithHeuristic struct {
	direction string
	heuristic int
}

func (sd sortedDirections) Len() int           { return len(sd) }
func (sd sortedDirections) Swap(i, j int)      { sd[i], sd[j] = sd[j], sd[i] }
func (sd sortedDirections) Less(i, j int) bool { return sd[i].heuristic < sd[j].heuristic }


func (s state) DebugPrint() {
	fmt.Println("Head: ", s.head)
	fmt.Println("Used directions", s.usedDirections)
	for y := 0;y < len(s.board[0]);y++ {
		for x := 0;x < len(s.board);x++ {
			fmt.Print(s.board[x][y], " ")
		}
		fmt.Println()
	}
}

func (s *state) Clone() *state {
	var newState state
	//Copy board
	board := make([][]int, len(s.board))
	for i := 0;i < len(s.board[0]);i++ {
		board[i] = make([]int, len(s.board[0]))
	}
	for x := 0;x < len(s.board);x++ {
		for y := 0;y < len(s.board[0]);y++ {
			board[x][y] = s.board[x][y]
		}
	}

	newState.head = s.head
	newState.board = board
	newState.usedDirections = make([]string, 0)
	newState.currentDirection = ""
	return &newState
}

func (s *state) isDirectionApplicable(d string) bool {
	x, y := 0, 0
	switch d {
	case "left":
		x = -1
	case "right":
		x = 1
	case "up":
		y = -1
	case "down":
		y = 1
	}
	if s.head.X + x < 0 {
		return false
	}
	if s.head.Y + y < 0 {
		return false
	}
	if s.head.X + x > len(s.board) - 1 {
		return false
	}
	if s.head.Y + y > len(s.board[0]) - 1 {
		return false
	}

	return s.board[s.head.X + x][s.head.Y + y] == 0
}

func (s *state) ApplyDirection(d string) {
	x, y := 0, 0
	switch d {
	case "left":
		x = -1
	case "right":
		x = 1
	case "up":
		y = -1
	case "down":
		y = 1
	}
	s.head.X += x
	s.head.Y += y
	s.board[s.head.X][s.head.Y] = 2
}

func (s *state) RevertDirection() {

}

func (bt *backTrack) heuristic(from, to Point, direction string) int {
	optimalDirection := ""
	if dX := to.X - from.X; dX > 0 {
		optimalDirection = "right"
	} else if dX < 0 {
		optimalDirection = "left"
	}
	if direction == optimalDirection {
		return 0
	}

	if dY := to.Y - from.Y; dY > 0 {
		optimalDirection = "down"
	} else if dY < 0 {
		optimalDirection = "up"
	}

	if optimalDirection == direction {
		return 0
	} else {
		return 1
	}
}

func (bt *backTrack) IsFinalState(s *state, from, to Point) bool {
	return from.X == to.X && from.Y == to.Y
}

func (bt *backTrack) FindPath(s state, from, to Point) []string {
	newState := s.Clone()

	//Find best direction which is applicable
	//Order directions by heuristic
	var sd sortedDirections
	sd = make([]directionWithHeuristic, 0)
	for _, d := range bt.directions {
		h := bt.heuristic(from,  to, d)
		sd = append(sd, directionWithHeuristic{d, h})
	}
	sort.Sort(sd)
	i := 0
	for len(newState.usedDirections) < 4 {
		newState.usedDirections = append(newState.usedDirections, sd[i].direction)
		if newState.isDirectionApplicable(sd[i].direction) {
			newState.ApplyDirection(sd[i].direction)
			if bt.IsFinalState(newState, newState.head, to) {
				break
			} else {
				path := bt.FindPath(*newState, newState.head, to)
				if path[0] != "" {
					return append(path, sd[i].direction)
				} else {
					//Back tracking
					newState.board[newState.head.X][newState.head.Y] = 0
					newState.head = s.head
				}
			}
		}
		i++
	}
	if i < 4 {
		return []string{sd[i].direction}
	} else {
		return []string{""}
	}

}

func GameFactory(data *MoveRequest) *Game {
	var head Point
	board := make([][]int, data.Width)
	for i := 0; i < data.Width; i++ {
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
	minDistance := math.Pow(float64(g.data.Width+g.data.Height), 2)

	for i, food := range g.data.Food {
		if d := math.Pow(float64(food.X-g.head.X), 2) + math.Pow(float64(food.Y-g.head.Y), 2); d < minDistance {
			minDistance = d
			minIndex = i
		}
	}

	return minIndex
}

func (g *Game) ChooseDirection() string {
	bt := backTrack{
		[]string{"up", "down", "left", "right"},
	}
	s := state{
		[]string{},
		"",
		g.board,
		g.head,
	}
	i := g.findClosestFood()
	start := time.Now()
	path := bt.FindPath(s, g.head, g.data.Food[i])
	fmt.Println(time.Since(start))

	return path[len(path) - 1]
}
