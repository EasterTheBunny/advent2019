package main

import (
	"testing"
)

func TestSwitchOps(t *testing.T) {
	values := []int{1, 0, 0, 0, 99}
	newValues := switchOps(0, values)

	if newValues[0] != 2 {
		t.Errorf("%v", values)
		t.Errorf("initial position value %v; expected %v", newValues[0], 2)
	}

	values = []int{2, 3, 0, 3, 99}
	newValues = switchOps(0, values)

	if newValues[3] != 6 {
		t.Errorf("%v", values)
		t.Errorf("position 3 value %v; expected %v", newValues[3], 6)
	}

	values = []int{2, 4, 4, 5, 99, 0}
	newValues = switchOps(0, values)

	if newValues[5] != 9801 {
		t.Errorf("%v", values)
		t.Errorf("position 5 value %v; expected %v", newValues[5], 9801)
	}

	values = []int{1, 1, 1, 4, 99, 5, 6, 0, 99}
	newValues = switchOps(0, values)

	if newValues[0] != 30 {
		t.Errorf("%v", values)
		t.Errorf("initial position value %v; expected %v", newValues[0], 30)
	}

	if newValues[4] != 2 {
		t.Errorf("%v", values)
		t.Errorf("position 4 value %v; expected %v", newValues[4], 2)
	}
}
