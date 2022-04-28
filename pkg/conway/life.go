package conway

import (
	"math/rand"
)

const (
	ALIVE = true
	DEAD  = false
)

// GameOfLife ...
type GameOfLife struct {
	size   int
	width  int
	height int
	state  []bool
	next   []bool
	offset []int
}

// NewGameOfLife creates a new game with width and height,
// it also serves as our initialization function 'setup()'
func NewGameOfLife(width, height int) *GameOfLife {
	game := new(GameOfLife)
	game.width = width
	game.height = height
	game.size = width * height
	game.state = make([]bool, game.size)
	game.next = make([]bool, game.size)
	game.offset = []int{
		-width - 1, // nw
		-width,     // n
		-width + 1, // ne
		1,          // e
		width + 1,  // se
		width,      // s
		width - 1,  // sw
		-1,         // w
	}
	game.Fill(25)
	return game
}

// Fills the state with dead cells and randomly seeds
// the state with live cells up to percentage
func (game *GameOfLife) Fill(percentage float64) {
	for i := range game.state {
		game.state[i] = DEAD
	}
	living := int(float64(game.size) * percentage / 100)
	for i := 0; i < living; i++ {
		game.state[rand.Intn(game.size)] = ALIVE
	}
}

// Seed sets a cell status to ALIVE
func (game *GameOfLife) Seed(x, y int) {
	game.state[y*game.width+x] = ALIVE
}

// Step creates the next generation of cells
func (game *GameOfLife) Step() {
	for i := range game.state {
		neighbours := 0
		for _, j := range game.offset {
			neighbours += game.At(i + j)
		}
		if game.state[i] == ALIVE && neighbours < 2 {
			game.next[i] = DEAD
		} else if game.state[i] == ALIVE && neighbours > 3 {
			game.next[i] = DEAD
		} else if game.state[i] == DEAD && neighbours == 3 {
			game.next[i] = ALIVE
		} else {
			game.next[i] = game.state[i]
		}
	}
	game.state, game.next = game.next, game.state
}

// At returns the cell 'status' at a given index (1D)
func (game *GameOfLife) At(i int) int {
	if i < 0 {
		i += game.size
	} else if i >= game.size {
		i -= game.size
	}
	if game.state[i] == ALIVE {
		return 1
	}
	return 0
}

// State returns the current generation's state
func (game *GameOfLife) State() []bool {
	return game.state
}

// Width returns the grid's width
func (game *GameOfLife) Width() int {
	return game.width
}

// Height returns the grid's height
func (game *GameOfLife) Height() int {
	return game.height
}
