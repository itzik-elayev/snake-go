package main

import (
	"errors"
	"math/rand"
	"slices"

	"github.com/gdamore/tcell"
)

var (
	FailedToCreateNewBoardErr     = errors.New("failed to create a new board")
	FailedToInitializeNewBoardErr = errors.New("failed to create a initialize new board")
)

type BoardLocation struct {
	x, y int
}

type Board struct {
	screen tcell.Screen
	snake  *Snake
	foods  []BoardLocation
}

func NewBoard() (*Board, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, errors.Join(FailedToCreateNewBoardErr, err)
	}

	if err = screen.Init(); err != nil {
		return nil, errors.Join(FailedToInitializeNewBoardErr, err)
	}

	screen.Clear()

	w, h := screen.Size()
	snake := NewSnake(BoardLocation{x: w / 2, y: h / 2})

	b := &Board{
		screen: screen,
		snake:  snake,
	}

	b.setBorders()

	return b, nil
}

func (b *Board) createFood() error {
	var (
		w, h                 = b.screen.Size()
		x, y, failedAttempts int
		maxFailedAttempts    = w * h
	)

	for {
		if failedAttempts == maxFailedAttempts {
			return errors.New("failed to find x,y coordinates to place food at (maybe you won the game?)")
		}

		x, y = 1+rand.Intn(w-2), 1+rand.Intn(h-2)

		if b.snake.isIn(x, y) || slices.Contains(b.foods, BoardLocation{x: x, y: y}) {
			failedAttempts++
			continue
		}

		break
	}

	b.foods = append(b.foods, BoardLocation{
		x: x,
		y: y,
	})

	return nil
}

func (b *Board) draw() {
	b.setBorders()

	for _, f := range b.foods {
		b.screen.SetContent(f.x, f.y, '■', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}

	for _, s := range b.snake.body {
		b.screen.SetContent(s.x, s.y, '■', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
	}

	b.screen.Show()
}

func (b *Board) setBorders() {
	var (
		border = '█'
	)

	w, h := b.screen.Size()

	for x := range w {
		b.screen.SetContent(x, 0, border, nil, tcell.StyleDefault)
		b.screen.SetContent(x, h-1, border, nil, tcell.StyleDefault)
	}

	for y := range h {
		b.screen.SetContent(0, y, border, nil, tcell.StyleDefault)
		b.screen.SetContent(w-1, y, border, nil, tcell.StyleDefault)
	}
}

func (b *Board) resize() {
	b.screen.Clear()
	b.setBorders()
}
