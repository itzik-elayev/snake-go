package main

import (
	"errors"
	"github.com/gdamore/tcell"
	"os"
	"slices"
	"time"
)

type GameState int

const (
	InitialGameState GameState = iota
	RunningGameState
	PausedGameState
	LostGameState
)

var (
	GameOverErr = errors.New("game over")
)

type Game struct {
	state    GameState
	board    *Board
	score    int
	numFoods int
}

func NewGame(board *Board, numFoods int) *Game {
	return &Game{
		state:    InitialGameState,
		board:    board,
		score:    0,
		numFoods: numFoods,
	}
}

func (g *Game) Start() error {
	var (
		err error
	)

	for i := 0; i < g.numFoods; i++ {
		if err = g.board.createFood(); err != nil {
			return err
		}
	}

	g.state = RunningGameState

	go g.handleInput()

	tick := time.Tick(33 * time.Millisecond)

	for {
		select {
		case <-tick:
			if err = g.update(); err != nil {
				return err
			}

			g.render()
		}
	}
}

func (g *Game) update() error {
	snakeTail := g.board.snake.tail()
	nextSnakeHead := g.board.snake.nextHead()

	w, h := g.board.screen.Size()

	if g.board.snake.isIn(nextSnakeHead.x, nextSnakeHead.y) || (nextSnakeHead.x <= 0 || nextSnakeHead.x >= w-1) || (nextSnakeHead.y <= 0 || nextSnakeHead.y >= h-1) {
		g.state = LostGameState
		return GameOverErr
	}

	g.board.snake.move(nextSnakeHead)

	foodIndex := slices.Index(g.board.foods, nextSnakeHead)
	if foodIndex != -1 {
		g.board.snake.grow(snakeTail)
		g.score++

		eatenFood := g.board.foods[foodIndex]
		g.board.screen.SetContent(eatenFood.x, eatenFood.y, ' ', nil, tcell.StyleDefault)
		g.board.foods = slices.Delete(g.board.foods, foodIndex, foodIndex+1)

		if err := g.board.createFood(); err != nil {
			return err
		}
	} else {
		g.board.screen.SetContent(snakeTail.x, snakeTail.y, ' ', nil, tcell.StyleDefault)
	}

	return nil
}

func (g *Game) render() {
	g.board.draw()
}

func (g *Game) handleInput() {
	for {
		switch ev := g.board.screen.PollEvent().(type) {
		case *tcell.EventResize:
			g.board.resize()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				g.board.screen.Fini()
				os.Exit(0)
			}

			if direction, ok := keyToDirection[ev.Key()]; ok {
				if g.board.snake.isValidDirection(direction) {
					g.board.snake.setDirection(direction)
				}
			}
		}
	}
}
