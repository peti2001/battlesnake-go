package main

import (
	"fmt"
	"math"
	"sort"
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
	counter    int
}

type sortedDirections []directionWithHeuristic

type directionWithHeuristic struct {
	direction string
	heuristic int
}

func (sd sortedDirections) Len() int           { return len(sd) }
func (sd sortedDirections) Swap(i, j int)      { sd[i], sd[j] = sd[j], sd[i] }
func (sd sortedDirections) Less(i, j int) bool { return sd[i].heuristic < sd[j].heuristic }

type sortedFood []foodWithDistance

type foodWithDistance struct {
	coord    Point
	distance float64
}

func (sf sortedFood) Len() int           { return len(sf) }
func (sf sortedFood) Swap(i, j int)      { sf[i], sf[j] = sf[j], sf[i] }
func (sf sortedFood) Less(i, j int) bool { return sf[i].distance < sf[j].distance }

func (s state) DebugPrint() {
	fmt.Println("Head: ", s.head)
	fmt.Println("Used directions", s.usedDirections)
	for y := 0; y < len(s.board[0]); y++ {
		for x := 0; x < len(s.board); x++ {
			fmt.Print(s.board[x][y], " ")
		}
		fmt.Println()
	}
}

func (s *state) Clone() *state {
	var newState state
	//Copy board
	board := make([][]int, len(s.board))
	for i := 0; i < len(s.board[0]); i++ {
		board[i] = make([]int, len(s.board[0]))
	}
	for x := 0; x < len(s.board); x++ {
		for y := 0; y < len(s.board[0]); y++ {
			board[x][y] = s.board[x][y]
		}
	}

	newState.head = s.head
	newState.board = board
	newState.usedDirections = make([]string, 0)
	newState.currentDirection = ""
	return &newState
}

func (s *state) IsDirectionApplicable(d string) bool {
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
	if s.head.X+x < 0 {
		return false
	}
	if s.head.Y+y < 0 {
		return false
	}
	if s.head.X+x > len(s.board)-1 {
		return false
	}
	if s.head.Y+y > len(s.board[0])-1 {
		return false
	}

	return s.board[s.head.X+x][s.head.Y+y] == 0 && s.HasEscapePath(d)
}

func (s *state) HasEscapePath(d string) bool {
	newState := s.Clone()
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
	newState.head.X += x
	newState.head.Y += y
	newState.board[newState.head.X][newState.head.Y] = 2
	var fn func(b [][]int, coord Point) int
	fn = func(b [][]int, coord Point) int {
		if coord.X < 0 || coord.Y < 0 {
			return 0
		}
		if coord.X >= len(b) || coord.Y >= len(b[0]) {
			return 0
		}
		if b[coord.X][coord.Y] == 0 {
			b[coord.X][coord.Y] = 2
			return 1 + fn(b, Point{coord.X + 1, coord.Y}) + fn(b, Point{coord.X, coord.Y + 1}) +
				fn(b, Point{coord.X - 1, coord.Y}) + fn(b, Point{coord.X, coord.Y - 1})
		} else {
			return 0
		}
	}
	count := fn(newState.board, Point{newState.head.X + 1, newState.head.Y}) + fn(newState.board, Point{newState.head.X, newState.head.Y + 1}) +
		fn(newState.board, Point{newState.head.X - 1, newState.head.Y}) + fn(newState.board, Point{newState.head.X, newState.head.Y - 1})
	lengthOfSnakes := 0
	for x := 0; x < len(s.board); x++ {
		for y := 0; y < len(s.board[0]); y++ {
			if s.board[x][y] == 1 {
				lengthOfSnakes++
			}
		}
	}
	return count > int(lengthOfSnakes/2)
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
	bt.counter++
	if bt.counter >= 50 {
		return []string{""}
	}
	newState := s.Clone()

	//Find best direction which is applicable
	//Order directions by heuristic
	var sd sortedDirections
	sd = make([]directionWithHeuristic, 0)
	for _, d := range bt.directions {
		h := bt.heuristic(from, to, d)
		sd = append(sd, directionWithHeuristic{d, h})
	}
	sort.Sort(sd)
	i := 0
	for len(newState.usedDirections) < 4 {
		newState.usedDirections = append(newState.usedDirections, sd[i].direction)
		if newState.IsDirectionApplicable(sd[i].direction) {
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

func (g *Game) findClosestFood() sortedFood {
	var sf sortedFood
	sf = make([]foodWithDistance, len(g.data.Food))
	for i, food := range g.data.Food {
		sf[i].coord = food
		sf[i].distance = math.Pow(float64(food.X-g.head.X), 2) + math.Pow(float64(food.Y-g.head.Y), 2)
	}
	sort.Sort(sf)

	return sf
}

func (g *Game) ChooseDirection() string {
	var path []string
	bt := backTrack{
		[]string{"up", "down", "left", "right"},
		0,
	}
	s := state{
		[]string{},
		"",
		g.board,
		g.head,
	}
	sf := g.findClosestFood()
	start := time.Now()
	for i := 0; i < len(g.data.Food); i++ {
		bt.counter = 0
		path = bt.FindPath(s, g.head, sf[i].coord)
		fmt.Println("PATH: ", path)
		if path[0] != "" {
			break
		}
	}
	fmt.Println(time.Since(start))

	return path[len(path)-1]
}
