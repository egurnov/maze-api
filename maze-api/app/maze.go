package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/egurnov/maze-api/maze-api/model"
	"github.com/egurnov/maze-api/maze-api/service"
	"github.com/gin-gonic/gin"
)

type MazeDTO struct {
	GridSize string   `json:"gridSize"`
	Entrance string   `json:"entrance"`
	Walls    []string `json:"walls"`
}

type CreateMazeDTO = MazeDTO
type MazeResponseDTO struct {
	ID int64 `json:"id"`
	MazeDTO
}

type GetAllMazesResponseDTO struct {
	Mazes []*MazeResponseDTO `json:"mazes"`
}

type SolutionResponseDTO struct {
	Path []string `json:"path"`
}

// CreateMaze godoc
// @Summary Create a new maze
// @ID CreateMaze
// @Tags Maze
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param maze body MazeDTO true "Maze description"
// @Success 201 {object} IDResponseDTO
// @Failure 400 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /maze [post]
func (a *App) CreateMaze(ctx *gin.Context) {
	var maze MazeDTO

	err := ctx.ShouldBindJSON(&maze)
	if err != nil {
		ctx.Error(err).SetType(BadRequestErrorType)
		return
	}

	rowcol := strings.Split(maze.GridSize, "x")
	if len(rowcol) != 2 {
		ctx.Error(errors.New("invalid grid size value")).SetType(BadRequestErrorType)
		return
	}
	rows, err := strconv.Atoi(rowcol[0])
	if err != nil {
		ctx.Error(err).SetType(BadRequestErrorType)
		return
	}
	cols, err := strconv.Atoi(rowcol[1])
	if err != nil {
		ctx.Error(err).SetType(BadRequestErrorType)
		return
	}
	if _, err := service.ParseCoords(maze.Entrance); err != nil {
		ctx.Error(errors.New("invalid entrance: " + maze.Entrance)).SetType(BadRequestErrorType)
		return
	}
	for _, wall := range maze.Walls {
		if wall == maze.Entrance {
			ctx.Error(errors.New("entrance cannot be a wall")).SetType(BadRequestErrorType)
			return
		}
		if _, err := service.ParseCoords(wall); err != nil {
			ctx.Error(errors.New("invalid wall: " + wall)).SetType(BadRequestErrorType)
			return
		}
	}

	id, err := a.MazeService.Create(&model.Maze{
		Rows:     rows,
		Cols:     cols,
		Entrance: maze.Entrance,
		Walls:    maze.Walls,
		UserID:   ctx.GetInt64(CTXUserID),
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, IDResponseDTO{ID: id})
}

// GetMaze godoc
// @Summary Get one specific maze belonging to the current user
// @ID GetMaze
// @Tags Maze
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param   id  path     integer     true  "maze id"
// @Success 200 {object} MazeResponseDTO
// @Failure 400 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /maze/{id} [get]
func (a *App) GetMaze(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 0, 64)
	if err != nil {
		ctx.Error(err).SetType(BadRequestErrorType)
		return
	}

	res, err := a.MazeService.GetByID(id, ctx.GetInt64(CTXUserID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, &MazeResponseDTO{
		ID: res.ID,
		MazeDTO: MazeDTO{
			GridSize: fmt.Sprintf("%dx%d", res.Rows, res.Cols),
			Entrance: res.Entrance,
			Walls:    res.Walls,
		},
	})
}

// GetAllMazes godoc
// @Summary Get all mazes belonging to the current user
// @ID GetAllMazes
// @Tags Maze
// @Accept json
// @Produce json
// @Security bearerAuth
// @Success 200 {object} GetAllMazesResponseDTO
// @Failure 400 {object} Message
// @Failure 500 {object} Message
// @Failure 500 {object} Message
// @Router /maze [get]
func (a *App) GetAllMazes(ctx *gin.Context) {
	res, err := a.MazeService.GetAll(ctx.GetInt64(CTXUserID))
	if err != nil {
		ctx.Error(err)
		return
	}

	allMazes := &GetAllMazesResponseDTO{Mazes: make([]*MazeResponseDTO, len(res))}
	for i, m := range res {
		allMazes.Mazes[i] = &MazeResponseDTO{
			ID: m.ID,
			MazeDTO: MazeDTO{
				GridSize: fmt.Sprintf("%dx%d", m.Rows, m.Cols),
				Entrance: m.Entrance,
				Walls:    m.Walls,
			},
		}
	}

	ctx.JSON(http.StatusOK, allMazes)
}

// SolveMaze godoc
// @Summary Solve a previously stored maze
// @ID SolveMaze
// @Tags Maze
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param   id  		path     integer    true  "maze id"
// @Param   steps   query     string     true  "Find shortest or longest path"       Enums(min, max)
// @Success 201 {object} SolutionResponseDTO
// @Failure 400 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /maze/{id}/solution [get]
func (a *App) SolveMaze(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 0, 64)
	if err != nil {
		ctx.Error(err).SetType(BadRequestErrorType)
		return
	}

	steps := ctx.Query("steps")
	if steps != "min" && steps != "max" {
		ctx.Error(errors.New("invalid steps value")).SetType(BadRequestErrorType)
		return
	}

	res, err := a.MazeService.Solve(id, ctx.GetInt64(CTXUserID), steps)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, &SolutionResponseDTO{
		Path: res,
	})
}
