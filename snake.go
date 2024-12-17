package main

import (
	"github.com/gdamore/tcell"
	"slices"
	"sync"
)

type SnakeDirection int

const (
	LeftSnakeDirection SnakeDirection = iota
	RightSnakeDirection
	UpSnakeDirection
	DownSnakeDirection
)

var (
	keyToDirection = map[tcell.Key]SnakeDirection{
		tcell.KeyLeft:  LeftSnakeDirection,
		tcell.KeyRight: RightSnakeDirection,
		tcell.KeyUp:    UpSnakeDirection,
		tcell.KeyDown:  DownSnakeDirection,
	}
)

type Snake struct {
	body           []BoardLocation
	direction      SnakeDirection
	directionMutex sync.Mutex
}

func NewSnake(location BoardLocation) *Snake {
	return &Snake{
		body:      []BoardLocation{location},
		direction: LeftSnakeDirection,
	}
}

func (s *Snake) move(location BoardLocation) {
	s.body = append(s.body, location)
	s.body = s.body[1:]
}

func (s *Snake) head() BoardLocation {
	return s.body[len(s.body)-1]
}

func (s *Snake) tail() BoardLocation {
	return s.body[0]
}

func (s *Snake) nextHead() BoardLocation {
	nextHead := s.head()

	direction := s.getDirection()

	switch direction {
	case LeftSnakeDirection:
		nextHead.x--
	case RightSnakeDirection:
		nextHead.x++
	case UpSnakeDirection:
		nextHead.y--
	case DownSnakeDirection:
		nextHead.y++
	}

	return nextHead
}

func (s *Snake) grow(next BoardLocation) {
	s.body = append([]BoardLocation{next}, s.body...)
}

func (s *Snake) isIn(x, y int) bool {
	return slices.Contains(s.body, BoardLocation{x: x, y: y})
}

func (s *Snake) getDirection() SnakeDirection {
	s.directionMutex.Lock()
	defer s.directionMutex.Unlock()
	return s.direction
}

func (s *Snake) setDirection(direction SnakeDirection) {
	s.directionMutex.Lock()
	s.direction = direction
	s.directionMutex.Unlock()
}

func (s *Snake) isValidDirection(direction SnakeDirection) bool {
	if len(s.body) == 1 {
		return true
	}

	currentDirection := s.getDirection()

	if (currentDirection == LeftSnakeDirection && direction == RightSnakeDirection) ||
		(currentDirection == RightSnakeDirection && direction == LeftSnakeDirection) ||
		(currentDirection == UpSnakeDirection && direction == DownSnakeDirection) ||
		(currentDirection == DownSnakeDirection && direction == UpSnakeDirection) {
		return false
	}

	return true
}
