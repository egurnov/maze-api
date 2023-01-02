package service

import (
	"errors"
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

func (s *MazeService) Solve(id, userId int64, steps string) ([]string, error) {
	maze, err := s.Store.GetByID(id, userId)
	if err != nil {
		return nil, err
	}

	return Solve(maze.Rows, maze.Cols, maze.Entrance, maze.Walls, steps)
}

func Solve(rows, cols int, entrance string, walls []string, steps string) ([]string, error) {
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

	// DEBUG
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
	//DEBUG

	start, err := ParseCoords(entrance)
	if err != nil {
		return nil, err
	}

	var res []Coords
	switch steps {
	case StepsMin:
		res, err = solveMin(maze, start)
	case StepsMax:
		err = errors.New("not implemented") //TODO
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

func solveMin(maze [][]bool, start Coords) ([]Coords, error) {
	type QueueEntry struct {
		Coords
		prev int
		len  int
	}
	q := []QueueEntry{{start, -1, 0}}
	rowLim := len(maze)
	colLim := len(maze[0])
	exit := -1

	been := make([][]bool, rowLim)
	for i := range been {
		been[i] = make([]bool, colLim)
	}
	been[start.Row][start.Col] = true

bfs:
	for i := 0; i < len(q); i++ {
		cur := q[i]

		// fmt.Printf(">>> Vising %v\n", cur)

		for _, delta := range []Coords{{+1, 0}, {-1, 0}, {0, +1}, {0, -1}} {

			// fmt.Printf(">>> Trying %v\n", Coords{cur.Row + delta.Row, cur.Col + delta.Col})

			if 0 <= cur.Row+delta.Row && cur.Row+delta.Row < rowLim &&
				0 <= cur.Col+delta.Col && cur.Col+delta.Col < colLim &&
				!maze[cur.Row+delta.Row][cur.Col+delta.Col] &&
				!been[cur.Row+delta.Row][cur.Col+delta.Col] {

				// fmt.Printf(">>> Adding %v\n", Coords{cur.Row + delta.Row, cur.Col + delta.Col})

				q = append(q, QueueEntry{
					Coords: Coords{cur.Row + delta.Row, cur.Col + delta.Col},
					prev:   i,
					len:    cur.len + 1,
				})
				been[cur.Row+delta.Row][cur.Col+delta.Col] = true

				if cur.Row+delta.Row == rowLim-1 {
					exit = len(q) - 1
					break bfs
				}
			}
		}
	}

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
