package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/egurnov/maze-api/maze-api/model"
)

type MazeService struct {
	Store model.MazeStore
}

var _ model.MazeService = &MazeService{}

func (s *MazeService) GetByID(id, userId int64) (*model.Maze, error) {
	return s.Store.GetByID(id, userId)
}

func (s *MazeService) GetAll(userId int64) ([]*model.Maze, error) {
	return s.Store.GetAll(userId)
}

func (s *MazeService) Create(maze *model.Maze) (int64, error) {
	return s.Store.Create(maze)
}

func (s *MazeService) Solve(ctx context.Context, id, userId int64, steps string) ([]string, error) {
	maze, err := s.Store.GetByID(id, userId)
	if err != nil {
		return nil, err
	}

	return Solve(ctx, maze.Rows, maze.Cols, maze.Entrance, maze.Walls, steps)
}

func Solve(ctx context.Context, rows, cols int, entrance string, walls []string, steps string) ([]string, error) {
	maze := make([][]bool, rows)
	for i := range maze {
		maze[i] = make([]bool, cols)
	}

	for _, wall := range walls {
		c, err := ParseCoords(wall)
		if err != nil {
			return nil, err
		}
		maze[c.Row][c.Col] = true
	}

	// printMaze(maze) // DEBUG

	start, err := ParseCoords(entrance)
	if err != nil {
		return nil, err
	}

	var res []Coords
	switch steps {
	case StepsMin:
		res, err = solveMin(ctx, maze, start)
	case StepsMax:
		res, err = solveMax(ctx, maze, start)
	default:
		return nil, model.ErrInvalidInput
	}
	if err != nil {
		return nil, err
	}

	resStr := make([]string, len(res))
	for i, c := range res {
		resStr[i] = CoordsToA1(c)
	}

	return resStr, nil
}

// SolveMax uses a non-recursive Depth First Search algorithm.
// For this kind of graph there is no polinomial time solution, so exponential is the best we can do.
// Because we manage our own stack, it also represents the current path.
func solveMax(ctx context.Context, maze [][]bool, start Coords) ([]Coords, error) {
	// Init
	type StackEntry struct {
		Coords
		deltaI int
	}
	st := []*StackEntry{{start, 0}}
	rows := len(maze)
	cols := len(maze[0])
	var res []Coords

	been := make([][]bool, rows)
	for i := range been {
		been[i] = make([]bool, cols)
	}
	been[start.Row][start.Col] = true

	deltas := []Coords{{+1, 0}, {-1, 0}, {0, +1}, {0, -1}}

	// Main DFS loop
dfs:
	for len(st) > 0 {
		// Timelimit check
		select {
		case <-ctx.Done():
			return nil, model.ErrorTimelimitReached
		default:
		}

		// Take stack top
		cur := st[len(st)-1]
		been[cur.Row][cur.Col] = true

		// fmt.Printf(">>> Vising %v\n", cur) // DEBUG

		// Check exit conditions
		if cur.Row == rows-1 {
			if len(st) > len(res) {
				res = make([]Coords, len(st))
				for i := range st {
					res[i] = st[i].Coords
				}
			}

			// Don't go anywhere from exit
			cur.deltaI = len(deltas)
		}

		// Iterate directions
		for cur.deltaI < len(deltas) {
			delta := deltas[cur.deltaI]
			cur.deltaI++

			// fmt.Printf(">>> Trying %v\n", Coords{cur.Row + delta.Row, cur.Col + delta.Col}) // DEBUG

			// If can go in this direction
			if 0 <= cur.Row+delta.Row && cur.Row+delta.Row < rows &&
				0 <= cur.Col+delta.Col && cur.Col+delta.Col < cols &&
				!maze[cur.Row+delta.Row][cur.Col+delta.Col] &&
				!been[cur.Row+delta.Row][cur.Col+delta.Col] {

				// fmt.Printf(">>> Adding %v\n", Coords{cur.Row + delta.Row, cur.Col + delta.Col}) // DEBUG

				// Add to stack
				st = append(st, &StackEntry{
					Coords: Coords{cur.Row + delta.Row, cur.Col + delta.Col},
				})
				continue dfs
			}
		}

		// Pop stack
		if cur.deltaI == len(deltas) {
			been[cur.Row][cur.Col] = false
			st = st[:len(st)-1]
		}
	}

	if res == nil {
		return nil, model.ErrorNoSolution
	}
	return res, nil
}

// SolveMin uses a non-recursive Breadth First Search algorithm. Visited cells are added to a queue and processed in order.
// Because there are no weights in the graph, all path lenghts in the queue will be in non-decreasing order.
// Execution time is O(number of reachable cells), in the worst case O(rows*columns).
func solveMin(ctx context.Context, maze [][]bool, start Coords) ([]Coords, error) {
	// Init
	type QueueEntry struct {
		Coords
		prev int // Index in the queue for the previous cell on the shortest path to this one
		len  int // Length of the sortest path to this cell
	}
	q := []QueueEntry{{start, -1, 0}}
	rows := len(maze)
	cols := len(maze[0])
	exit := -1 // Index in the queue for the exit cell

	been := make([][]bool, rows)
	for i := range been {
		been[i] = make([]bool, cols)
	}
	been[start.Row][start.Col] = true

	// Main BFS loop
bfs:
	for i := 0; i < len(q); i++ {
		// Timelimit check
		select {
		case <-ctx.Done():
			return nil, model.ErrorTimelimitReached
		default:
		}

		// Take the next cell from the queue
		cur := q[i]

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

				// Add to queue
				q = append(q, QueueEntry{
					Coords: Coords{cur.Row + delta.Row, cur.Col + delta.Col},
					prev:   i,
					len:    cur.len + 1,
				})
				// Marking visited cells early to avoid adding them multiple times
				been[cur.Row+delta.Row][cur.Col+delta.Col] = true

				// Check exit condition. The first path we find is the shortest.
				if cur.Row+delta.Row == rows-1 {
					exit = len(q) - 1
					break bfs
				}
			}
		}
	}

	// Build path
	if exit > 0 {
		res := make([]Coords, q[exit].len+1)
		for i := exit; i >= 0; i = q[i].prev {
			res[q[i].len] = q[i].Coords
		}
		return res, nil
	}

	return nil, model.ErrorNoSolution
}

// 0-based coordinates in the maze. Valid values are 0 <= row < rows, 0 <= col < cols.
type Coords struct {
	Row, Col int
}

const (
	StepsMin = "min"
	StepsMax = "max"
)

func ParseCoords(s string) (Coords, error) {
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

func AreValid(c Coords, rows, cols int) bool {
	return 0 <= c.Row && c.Row < rows && 0 <= c.Col && c.Col < cols
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

func printMaze(maze [][]bool) {
	for _, row := range maze {
		for _, cell := range row {
			if cell {
				fmt.Print(" ")
			} else {
				fmt.Print("X")
			}
		}
		fmt.Println()
	}
}
