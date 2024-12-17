package main

func main() {
	board, err := NewBoard()
	if err != nil {
		panic(err)
	}

	game := NewGame(board, 5)

	if err := game.Start(); err != nil {
		panic(err)
	}
}
