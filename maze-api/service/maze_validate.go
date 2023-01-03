package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/egurnov/maze-api/maze-api/model"
)

// 0-based coordinates in the maze. Valid values are 0 <= row < rows, 0 <= col < cols.
type Coords struct {
	Row, Col int
}

const (
	StepsMin = "min"
	StepsMax = "max"
)

func A1ToCoords(s string) (Coords, error) {
	row := 0
	col := 0
	i := 0
	for ; i < len(s) && 'A' <= s[i] && s[i] <= 'Z'; i++ {
		col = int(s[i]-'A'+1) + col*int('Z'-'A'+1)
	}
	if i == 0 {
		return Coords{}, model.ErrInvalidInput
	}
	for ; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		row = int(s[i]-'0') + row*int('9'-'0'+1)
	}
	if i != len(s) {
		return Coords{}, model.ErrInvalidInput
	}
	return Coords{row - 1, col - 1}, nil
}

func CoordsToA1(coords Coords) string {
	res := ""
	col := coords.Col + 1
	for ; col > 0; col /= int('Z' - 'A' + 1) {
		if col%int('Z'-'A'+1) == 0 {
			res = "Z" + res
			col -= 26
		} else {
			res = string(byte('A'+col%int('Z'-'A'+1)-1)) + res
		}

	}
	return res + strconv.Itoa(coords.Row+1)
}

func areValid(c Coords, rows, cols int) bool {
	return 0 <= c.Row && c.Row < rows && 0 <= c.Col && c.Col < cols
}

func parseGridSize(gridSize string) (rows, cols int, err error) {
	rowcol := strings.Split(gridSize, "x")
	if len(rowcol) != 2 {
		return 0, 0, errors.New("invalid grid size value")
	}

	rows, err = strconv.Atoi(rowcol[0])
	if err != nil {
		return 0, 0, errors.New("invalid rows value: " + err.Error())
	}

	cols, err = strconv.Atoi(rowcol[1])
	if err != nil {
		return 0, 0, errors.New("invalid cols value: " + err.Error())
	}

	return rows, cols, err
}

// true = wall, false = open
func makeMaze(rows, cols int, walls []string) ([][]bool, error) {
	maze := make([][]bool, rows)
	for i := range maze {
		maze[i] = make([]bool, cols)
	}

	for _, wall := range walls {
		c, err := A1ToCoords(wall)
		if err != nil {
			return nil, err
		}
		maze[c.Row][c.Col] = true
	}

	return maze, nil
}

//nolint:unused
func printMaze(maze [][]bool) {
	fPrintMaze(maze, os.Stdout)
}

func fPrintMaze(maze [][]bool, f io.Writer) {
	for _, row := range maze {
		for _, cell := range row {
			if cell {
				fmt.Fprint(f, "|X")
			} else {
				fmt.Fprint(f, "|_")
			}
		}
		fmt.Fprintln(f, "|")
	}
}

func ValidateMaze(gridSize, entrance string, walls []string) (rows, cols int, entranceCoords Coords, wallsCoords []Coords, err error) {
	// Grid size validation
	rows, cols, err = parseGridSize(gridSize)
	if err != nil {
		return 0, 0, Coords{}, nil, err
	}

	// Entrance validation
	entranceCoords, err = A1ToCoords(entrance)
	if err != nil || !areValid(entranceCoords, rows, cols) {
		return 0, 0, Coords{}, nil, errors.New("invalid entrance: " + entrance)
	}

	// Walls validation
	for _, wall := range walls {
		if wall == entrance {
			return 0, 0, Coords{}, nil, errors.New("entrance cannot be a wall")
		}

		coords, err := A1ToCoords(wall)
		if err != nil || !areValid(coords, rows, cols) {
			return 0, 0, Coords{}, nil, errors.New("invalid wall: " + wall)
		}

		wallsCoords = append(wallsCoords, coords)
	}

	// Exit point validation
	maze, err := makeMaze(rows, cols, walls)
	if err != nil {
		return 0, 0, Coords{}, nil, err
	}

	// TODO: If there are multiple open cells in the last row, but only one of them is directly accessible, is this a valid maze?
	count, err := countReachableExits(maze, entranceCoords)
	if err != nil {
		return 0, 0, Coords{}, nil, err
	}

	if count != 1 {
		return 0, 0, Coords{}, nil, errors.New("invalid exit point")
	}

	return rows, cols, entranceCoords, wallsCoords, nil
}

// countReachableExits uses a non-recursive Breadth First Search algorithm. Visited cells are added to a queue and processed in order.
// Execution time is O(number of reachable cells), in the worst case O(rows*columns).
func countReachableExits(maze [][]bool, start Coords) (int, error) {
	// Init
	type QueueEntry = Coords
	q := []QueueEntry{start}
	rows := len(maze)
	cols := len(maze[0])
	exitCount := 0

	been := make([][]bool, rows)
	for i := range been {
		been[i] = make([]bool, cols)
	}
	been[start.Row][start.Col] = true

	// Main BFS loop
	for i := 0; i < len(q); i++ {
		// Take the next cell from the queue
		cur := q[i]

		// Check exit condition
		if cur.Row == rows-1 {
			exitCount++
		}

		// fmt.Printf(">>> Vising %v\n", cur) // DEBUG

		// Iterate directions
		for _, delta := range []Coords{{+1, 0}, {-1, 0}, {0, +1}, {0, -1}} {

			// fmt.Printf(">>> Trying %v\n", Coords{cur.Row + delta.Row, cur.Col + delta.Col}) // DEBUG

			// If can go in this direction
			if 0 <= cur.Row+delta.Row && cur.Row+delta.Row < rows &&
				0 <= cur.Col+delta.Col && cur.Col+delta.Col < cols &&
				!maze[cur.Row+delta.Row][cur.Col+delta.Col] &&
				!been[cur.Row+delta.Row][cur.Col+delta.Col] {

				// fmt.Printf(">>> Adding %v\n", Coords{cur.Row + delta.Row, cur.Col + delta.Col}) // DEBUG

				// Marking visited cells early to avoid adding them multiple times
				been[cur.Row+delta.Row][cur.Col+delta.Col] = true

				// Add to queue
				q = append(q, Coords{cur.Row + delta.Row, cur.Col + delta.Col})
			}
		}
	}

	return exitCount, nil
}
