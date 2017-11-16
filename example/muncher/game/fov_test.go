package game

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
		expectedPos := Position{XPos: x, YPos: 0}
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
		expectedPos := Position{XPos: 5, YPos: y}
		if testPos != expectedPos {
			t.Errorf("Expected %+v but got %+v", expectedPos, testPos)
		}
	}
}

func TestDiagonalDownRight(t *testing.T) {
	cells := PlotLine(0, 1, 6, 4)

	expected := []Position{
		Position{XPos: 0, YPos: 1},
		Position{XPos: 1, YPos: 1},
		Position{XPos: 2, YPos: 2},
		Position{XPos: 3, YPos: 2},
		Position{XPos: 4, YPos: 3},
		Position{XPos: 5, YPos: 3},
		Position{XPos: 6, YPos: 4},
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
