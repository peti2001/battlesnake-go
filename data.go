package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type GameStartRequest struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	GameId string `json:"game_id"`
}

type GameStartResponse struct {
	Color   string  `json:"color"`
	HeadUrl *string `json:"head_url,omitempty"`
	Name    string  `json:"name"`
	Taunt   *string `json:"taunt,omitempty"`
}

type MoveRequest struct {
	You    string  `json:"you"`
	Food   []Point `json:"food"`
	GameId string  `json:"game_id"`
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Turn   int     `json:"turn"`
	Snakes []Snake `json:"snakes"`
}

type MoveResponse struct {
	Move  string  `json:"move"`
	Taunt *string `json:"taunt,omitempty"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Coords       []Point `json:"coords"`
	HealthPoints int     `json:"health_points"`
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Taunt        string  `json:"taunt"`
}

func NewMoveRequest(req *http.Request) (*MoveRequest, error) {
	body, _ := ioutil.ReadAll(req.Body)
	decoded := MoveRequest{}
	err := json.Unmarshal(body, &decoded)
	return &decoded, err
}

func NewGameStartRequest(req *http.Request) (*GameStartRequest, error) {
	decoded := GameStartRequest{}
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &decoded)

	return &decoded, err
}

func (snake Snake) Head() Point { return snake.Coords[0] }

// Decode [number, number] JSON array into a Point
func (point *Point) UnmarshalJSON(data []byte) error {
	var coords []int
	json.Unmarshal(data, &coords)
	if len(coords) != 2 {
		return errors.New("Bad set of coordinates: " + string(data))
	}
	*point = Point{X: coords[0], Y: coords[1]}
	return nil
}
