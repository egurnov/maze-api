package service_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/egurnov/maze-api/maze-api/service"
)

func TestParseCoords(t *testing.T) {
	testCases := []struct {
		in     string
		expRow int
		expCol int
		expErr bool
	}{
		{"A1", 0, 0, false},
		{"A10", 9, 0, false},
		{"A11", 10, 0, false},
		{"A100", 99, 0, false},
		{"A123", 122, 0, false},
		{"B2", 1, 1, false},
		{"Z1", 0, 25, false},
		{"AA1", 0, 26, false},
		{"AZ1", 0, 51, false},
		{"BA1", 0, 52, false},
		{"ZZ1", 0, 701, false},
		{"AAA1", 0, 702, false},
		{"a1", 0, 0, true},
		{"1a", 0, 0, true},
	}
	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			g := NewWithT(t)

			coords, err := service.ParseCoords(tc.in)

			g.Expect(err != nil).To(Equal(tc.expErr))
			g.Expect(coords).To(Equal(service.Coords{tc.expRow, tc.expCol}))
		})
	}
}
func TestCoordsToA1(t *testing.T) {
	testCases := []struct {
		exp string
		row int
		col int
	}{
		{"A1", 0, 0},
		{"A10", 9, 0},
		{"A11", 10, 0},
		{"A100", 99, 0},
		{"A123", 122, 0},
		{"B2", 1, 1},
		{"Z1", 0, 25},
		{"AA1", 0, 26},
		{"AZ1", 0, 51},
		{"BA1", 0, 52},
		{"ZZ1", 0, 701},
		{"AAA1", 0, 702},
	}
	for _, tc := range testCases {
		t.Run(tc.exp, func(t *testing.T) {
			g := NewWithT(t)

			res := service.CoordsToA1(service.Coords{tc.row, tc.col})

			g.Expect(res).To(Equal(tc.exp))
		})
	}
}

func TestSolveMin(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		ctx := context.Background()

		res, err := service.Solve(ctx, 8, 8, "A1", []string{"C1", "G1", "A2", "C2", "E2", "G2", "C3", "E3", "B4", "C4", "E4", "F4", "G4", "B5", "E5", "B6", "D6", "E6", "G6", "H6", "B7", "D7", "G7", "B8"}, "min")

		g := NewWithT(t)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(res).To(Equal([]string{"A1", "B1", "B2", "B3", "A3", "A4", "A5", "A6", "A7", "A8"}))
	})

	t.Run("no solution", func(t *testing.T) {
		ctx := context.Background()

		res, err := service.Solve(ctx, 8, 8, "A1", []string{"C1", "G1", "A2", "C2", "E2", "G2", "C3", "E3", "B4", "C4", "E4", "F4", "G4", "B5", "E5", "B6", "D6", "E6", "G6", "H6", "B7", "D7", "G7", "B8", "A8"}, "min")

		g := NewWithT(t)
		g.Expect(err).To(MatchError("no solution"))
		g.Expect(res).To(BeNil())
	})

	t.Run("2 paths", func(t *testing.T) {
		ctx := context.Background()

		res, err := service.Solve(ctx, 8, 8, "A1", []string{"C1", "G1", "A2", "C2", "D2", "E2", "G2", "E3", "B4", "C4", "E4", "F4", "G4", "C6", "E5", "B6", "D6", "E6", "G6", "H6", "B7", "D7", "G7", "B8"}, "min")

		g := NewWithT(t)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(res).To(Equal([]string{"A1", "B1", "B2", "B3", "A3", "A4", "A5", "A6", "A7", "A8"}))
	})
}

func TestSolveMax(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		ctx := context.Background()

		res, err := service.Solve(ctx, 8, 8, "A1", []string{"C1", "G1", "A2", "C2", "E2", "G2", "C3", "E3", "B4", "C4", "E4", "F4", "G4", "B5", "E5", "B6", "D6", "E6", "G6", "H6", "B7", "D7", "G7", "B8"}, "max")

		g := NewWithT(t)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(res).To(Equal([]string{"A1", "B1", "B2", "B3", "A3", "A4", "A5", "A6", "A7", "A8"}))
	})

	t.Run("no solution", func(t *testing.T) {
		ctx := context.Background()

		res, err := service.Solve(ctx, 8, 8, "A1", []string{"C1", "G1", "A2", "C2", "E2", "G2", "C3", "E3", "B4", "C4", "E4", "F4", "G4", "B5", "E5", "B6", "D6", "E6", "G6", "H6", "B7", "D7", "G7", "B8", "A8"}, "max")

		g := NewWithT(t)
		g.Expect(err).To(MatchError("no solution"))
		g.Expect(res).To(BeNil())
	})

	t.Run("2 paths", func(t *testing.T) {
		ctx := context.Background()

		res, err := service.Solve(ctx, 8, 8, "A1", []string{"C1", "G1", "A2", "C2", "D2", "E2", "G2", "E3", "B4", "C4", "E4", "F4", "G4", "C6", "E5", "B6", "D6", "E6", "G6", "H6", "B7", "D7", "G7", "B8"}, "max")

		g := NewWithT(t)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(res).To(Equal([]string{"A1", "B1", "B2", "B3", "C3", "D3", "D4", "D5", "C5", "B5", "A5", "A6", "A7", "A8"}))
	})

	t.Run("timelimit", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		res, err := service.Solve(ctx, 10, 10, "A1", []string{"A10", "B10", "C10", "D10", "E10", "F10", "G10", "H10", "I10"}, "max")

		g := NewWithT(t)
		g.Expect(err).To(MatchError("time limit reached"))
		g.Expect(res).To(BeNil())
	})
}
