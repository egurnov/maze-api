package store

import (
	"strings"

	"github.com/egurnov/maze-api/maze-api/model"
)

var _ model.MazeStore = &MazeStore{&Store{}}

type MazeStore struct{ *Store }

type Maze struct {
	ID       int64  `gorm:"primary_key;auto_increment"`
	Rows     int    `gorm:"not null;type:int"`
	Cols     int    `gorm:"not null;type:int"`
	Entrance string `gorm:"not null;type:varchar(100)"`
	Walls    string `gorm:"not null;type:varchar(500)"`

	UserID int64 `json:"-"`
}

func (s *MazeStore) GetByID(id, userId int64) (*model.Maze, error) {
	var maze Maze
	err := s.db.Where("user_id = ?", userId).First(&maze, id).Error
	return &model.Maze{
		ID:       maze.ID,
		Rows:     maze.Rows,
		Cols:     maze.Cols,
		Entrance: maze.Entrance,
		Walls:    strings.Split(maze.Walls, ","),
		UserID:   maze.UserID,
	}, wrapError(err)
}

func (s *MazeStore) GetAll(userId int64) ([]*model.Maze, error) {
	var mazes []*Maze
	err := s.db.Where("user_id = ?", userId).Find(&mazes).Error
	res := make([]*model.Maze, len(mazes))
	for i, maze := range mazes {
		res[i] = &model.Maze{
			ID:       maze.ID,
			Rows:     maze.Rows,
			Cols:     maze.Cols,
			Entrance: maze.Entrance,
			Walls:    strings.Split(maze.Walls, ","),
			UserID:   maze.UserID,
		}
	}
	return res, wrapError(err)
}

func (s *MazeStore) Create(maze *model.Maze) (int64, error) {
	dbMaze := Maze{
		Rows:     maze.Rows,
		Cols:     maze.Cols,
		Entrance: maze.Entrance,
		Walls:    strings.Join(maze.Walls, ","),
		UserID:   maze.UserID,
	}
	err := s.db.Create(&dbMaze).Error

	return dbMaze.ID, wrapError(err)
}
