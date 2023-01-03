package service

import (
	"bytes"
	"context"

	"github.com/egurnov/maze-api/maze-api/model"
)

type MazeService struct {
	Store model.MazeStore
}

var _ model.MazeService = &MazeService{}

func (s *MazeService) GetByID(id, userId int64) (*model.Maze, error) {
	return s.Store.GetByID(id, userId)
}

func (s *MazeService) PrintMaze(id, userId int64) ([]byte, error) {
	mazeDescr, err := s.Store.GetByID(id, userId)
	if err != nil {
		return nil, err
	}

	maze, err := makeMaze(mazeDescr.Rows, mazeDescr.Cols, mazeDescr.Walls)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	fPrintMaze(maze, b)

	return b.Bytes(), nil
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
