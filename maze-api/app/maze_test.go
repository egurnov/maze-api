package app_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/egurnov/maze-api/maze-api/app"
)

func TestValidateMaze(t *testing.T) {
	testCases := []struct {
		desc   string
		maze   app.MazeDTO
		expErr string
	}{
		{
			desc: "example",
			maze: app.MazeDTO{
				GridSize: "8x8",
				Entrance: "A1",
				Walls:    []string{"C1", "G1", "A2", "C2", "E2", "G2", "C3", "E3", "B4", "C4", "E4", "F4", "G4", "B5", "E5", "B6", "D6", "E6", "G6", "H6", "B7", "D7", "G7", "B8", "A8"},
			},
			expErr: "",
		},
		{
			desc: "valid",
			maze: app.MazeDTO{
				GridSize: "4x4",
				Entrance: "A1",
				Walls:    []string{"A4", "B4", "C4"},
			},
			expErr: "",
		},
		{
			desc: "invalid grid",
			maze: app.MazeDTO{
				GridSize: "4-4",
				Entrance: "A1",
				Walls:    []string{"A4", "B4", "C4"},
			},
			expErr: "invalid grid size",
		},
		{
			desc: "invalid rows",
			maze: app.MazeDTO{
				GridSize: "x4",
				Entrance: "A1",
				Walls:    []string{"A4", "B4", "C4"},
			},
			expErr: "invalid rows",
		},
		{
			desc: "invalid cols",
			maze: app.MazeDTO{
				GridSize: "4x",
				Entrance: "A1",
				Walls:    []string{"A4", "B4", "C4"},
			},
			expErr: "invalid cols",
		},
		{
			desc: "entrance out of bounds",
			maze: app.MazeDTO{
				GridSize: "4x4",
				Entrance: "E1",
				Walls:    []string{"A4", "B4", "C4"},
			},
			expErr: "invalid entrance",
		},
		{
			desc: "invalid entrance format",
			maze: app.MazeDTO{
				GridSize: "4x4",
				Entrance: "1A",
				Walls:    []string{"A4", "B4", "C4"},
			},
			expErr: "invalid entrance",
		},
		{
			desc: "entrance is a wall",
			maze: app.MazeDTO{
				GridSize: "4x4",
				Entrance: "A1",
				Walls:    []string{"A4", "B4", "C4", "A1"},
			},
			expErr: "entrance cannot be a wall",
		},
		{
			desc: "invalid wall format",
			maze: app.MazeDTO{
				GridSize: "4x4",
				Entrance: "A1",
				Walls:    []string{"A4", "B4", "C4", "1A"},
			},
			expErr: "invalid wall",
		},
		{
			desc: "no exit",
			maze: app.MazeDTO{
				GridSize: "4x4",
				Entrance: "A1",
				Walls:    []string{"A4", "B4", "C4", "D4"},
			},
			expErr: "invalid exit",
		},
		{
			desc: "2 exits",
			maze: app.MazeDTO{
				GridSize: "4x4",
				Entrance: "A1",
				Walls:    []string{"A4", "B4"},
			},
			expErr: "invalid exit",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			_, _, err := app.ValidateMaze(tc.maze)
			g := NewWithT(t)
			if len(tc.expErr) == 0 {
				g.Expect(err).ToNot(HaveOccurred())
			} else {
				g.Expect(err).To(MatchError(ContainSubstring(tc.expErr)))
			}
		})
	}
}
