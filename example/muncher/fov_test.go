package main

import (
	"testing"
)

func TestPlotLineHorizontalRight(t *testing.T) {
	cells := PlotLine(0, 0, 5, 0)

	if len(cells) != 6 {
		t.Fatalf("Expected result array length of 6, got %v instead", len(cells))
	}
	for x := 0; x <= 5; x++ {
		testPos := cells[x]
		expectedPos := Position{X: x, Y: 0}
		if testPos != expectedPos {
			t.Errorf("Expected %+v but got %+v", expectedPos, testPos)
		}
	}
}

func TestPLotLineVertialDown(t *testing.T) {
	cells := PlotLine(5, 5, 5, 10)

	if len(cells) != 6 {
		t.Fatalf("Expected result array length of 6, got %v instead", len(cells))
	}
	for y := 5; y <= 10; y++ {
		testPos := cells[y-5]
		expectedPos := Position{X: 5, Y: y}
		if testPos != expectedPos {
			t.Errorf("Expected %+v but got %+v", expectedPos, testPos)
		}
	}
}

func TestDiagonalDownRight(t *testing.T) {
	cells := PlotLine(0, 1, 6, 4)

	expected := []Position{
		Position{X: 0, Y: 1},
		Position{X: 1, Y: 1},
		Position{X: 2, Y: 2},
		Position{X: 3, Y: 2},
		Position{X: 4, Y: 3},
		Position{X: 5, Y: 3},
		Position{X: 6, Y: 4},
	}

	if len(cells) != len(expected) {
		t.Fatalf("Expected result array length of %v, got %v instead", len(expected), len(cells))
	}

	for index, pos := range cells {
		if pos != expected[index] {
			t.Errorf("At index %v got %+v but expected %+v", index, pos, expected)
		}
	}
}
