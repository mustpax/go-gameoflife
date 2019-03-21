package main

import "fmt"

type game struct {
	cells [][]chan bool
	rows  int
	cols  int
	done  chan bool
}

func newGame(rows int, cols int) *game {
	cells := make([][]chan bool, rows)
	for i := 0; i < rows; i++ {
		cells[i] = make([]chan bool, cols)
		for j := 0; j < cols; j++ {
			cells[i][j] = make(chan bool, 8)
		}
	}

	return &game{
		cells: cells,
		rows:  rows,
		cols:  cols,
		done:  make(chan bool),
	}
}

func (g *game) neighbors(row int, col int) []chan bool {
	ret := make([]chan bool, 0)

	if row > 0 {
		ret = append(ret, g.cells[row-1][col])
		if col > 0 {
			ret = append(ret, g.cells[row-1][col-1])
		}
		if col < (g.cols - 1) {
			ret = append(ret, g.cells[row-1][col+1])
		}
	}

	if row < (g.rows - 1) {
		ret = append(ret, g.cells[row+1][col])
		if col > 0 {
			ret = append(ret, g.cells[row+1][col-1])
		}
		if col < (g.cols - 1) {
			ret = append(ret, g.cells[row+1][col+1])
		}
	}

	if col > 0 {
		ret = append(ret, g.cells[row][col-1])
	}

	if col < (g.cols - 1) {
		ret = append(ret, g.cells[row][col+1])
	}
	return ret
}

func (g *game) startAgents(initialState [][]bool, generations int) {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			g.startAgent(initialState[i][j], generations, i, j)
		}
	}
}

func (g *game) startAgent(alive bool, generations int, row int, col int) {
	go func() {
		fmt.Printf("START %v:%v %v\n", row, col, alive)
		neighbs := g.neighbors(row, col)
		self := g.cells[row][col]
		for i := 0; i < generations; i++ {
			for _, n := range neighbs {
				n <- alive
			}
			neighbsAlive := 0
			for _ = range neighbs {
				if <-self {
					neighbsAlive++
				}
			}
			fmt.Printf("G%04d %v:%v %v %v\n", i, row, col, alive, neighbsAlive)
			if neighbsAlive < 2 {
				alive = false
			} else if neighbsAlive == 2 {
				// do nothing
			} else if neighbsAlive == 3 {
				alive = true
			} else {
				alive = false
			}
		}
		fmt.Printf("END   %v:%v %v\n", row, col, alive)
		g.done <- alive
	}()
}

func (g *game) waitForEnd() int {
	ret := 0
	cells := g.rows * g.cols
	fmt.Println("cells", cells)
	for i := 0; i < cells; i++ {
		if <-g.done {
			ret++
		}
	}
	return ret
}

func main() {
	g := newGame(2, 2)
	g.startAgents([][]bool{
		[]bool{true, true},
		[]bool{true, false},
	}, 1)
	totalAlive := g.waitForEnd()
	fmt.Println("Game is over, total alive:", totalAlive)
}
