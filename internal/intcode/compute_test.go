package intcode

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestSplitOp(t *testing.T) {
	inputs := []int{1002, 11101, 3, 4, 105, 1006, 7, 8, 99}

	expectedOps := []OpCode{
		MultiplyOp,
		AddOp,
		InputOp,
		OutputOp,
		JumpTrue,
		JumpFalse,
		LessThan,
		Equals,
		TerminateOp}

	expectedMod := [][]ParameterMode{
		[]ParameterMode{PositionMode, ImmediateMode, PositionMode, PositionMode},
		[]ParameterMode{ImmediateMode, ImmediateMode, ImmediateMode, PositionMode},
		[]ParameterMode{PositionMode, PositionMode, PositionMode, PositionMode},
		[]ParameterMode{PositionMode, PositionMode, PositionMode, PositionMode},
		[]ParameterMode{ImmediateMode, PositionMode, PositionMode, PositionMode},
		[]ParameterMode{PositionMode, ImmediateMode, PositionMode, PositionMode},
		[]ParameterMode{PositionMode, PositionMode, PositionMode, PositionMode},
		[]ParameterMode{PositionMode, PositionMode, PositionMode, PositionMode},
		[]ParameterMode{PositionMode, PositionMode, PositionMode, PositionMode}}

	for i, input := range inputs {
		code, modes := SplitOp(input)

		if code != expectedOps[i] {
			t.Errorf("incorrect op code %v for input %v; expected %v", code, input, expectedOps[i])
		}

		if len(modes) != len(expectedMod[i]) {
			t.Errorf("incorrect parameter length %v; expected %v", len(modes), len(expectedMod[i]))
		} else {
			for j, mod := range expectedMod[i] {
				if modes[j] != mod {
					t.Errorf("incorrect mode value %v; expected %v", modes[j], mod)
				}
			}
		}
	}
}

func TestGetParameters(t *testing.T) {
	inputs := [][]int{
		[]int{2, 3, 1, 1, 99},
		[]int{3, 4, 1, 4, 3, 3, 99},
		[]int{5, 6, 4, 1, 4, 3, 3, 99},
		[]int{3, 0, 4, 0, 99}}

	index := []int{1, 3, 1, 3}
	length := []int{3, 3, 2, 1}
	expected := [][]int{
		[]int{3, 1, 1},
		[]int{4, 3, 3},
		[]int{6, 4},
		[]int{0}}

	for i, input := range inputs {
		p := GetParameters(length[i], index[i], input)

		for j, b := range expected[i] {
			if p[j].Value != b {
				t.Errorf("incorrect parameter %v; expected %v", p[j].Value, b)
			}
		}
	}
}

func TestGetInstruction_Modes(t *testing.T) {
	inputs := [][]int{
		[]int{1002, 0, 2, 0, 99},
		[]int{3, 3, 11101, 3, 4, 0, 99},
		[]int{3, 0, 4, 0, 99},
		[]int{3, 3, 1105, 4, 9, 1101, 0, 0, 12, 4, 12, 99, 1}}

	inputInd := []int{0, 2, 0, 2}
	expectedMod := [][]ParameterMode{
		[]ParameterMode{PositionMode, ImmediateMode, PositionMode, PositionMode},
		[]ParameterMode{ImmediateMode, ImmediateMode, ImmediateMode, PositionMode},
		[]ParameterMode{PositionMode, PositionMode, PositionMode, PositionMode},
		[]ParameterMode{ImmediateMode, ImmediateMode, PositionMode, PositionMode}}

	for i, input := range inputs {
		inst := GetInstruction(inputInd[i], input)

		if len(inst.Modes) != len(expectedMod[i]) {
			t.Errorf("incorrect parameter length %v; expected %v", len(inst.Modes), len(expectedMod[i]))
		} else {
			for j, mod := range expectedMod[i] {
				if inst.Modes[j] != mod {
					t.Errorf("incorrect mode value %v; expected %v", inst.Modes[j], mod)
				}
			}
		}
	}
}

func TestGetInstruction_Parameter(t *testing.T) {
	inputs := [][]int{
		[]int{1002, 0, 2, 0, 99},
		[]int{3, 3, 11101, 3, 4, 0, 99},
		[]int{3, 0, 99},
		[]int{4, 0, 99},
		[]int{5, 0, 2, 99},
		[]int{6, 0, 2, 99},
		[]int{7, 0, 2, 99},
		[]int{8, 0, 2, 99}}

	inputInd := []int{0, 2, 0, 0, 0, 0, 0, 0}
	expected := []int{3, 3, 1, 1, 2, 2, 3, 3}

	for i, input := range inputs {
		inst := GetInstruction(inputInd[i], input)

		if len(inst.Parameters) != expected[i] {
			t.Errorf("incorrect number of parameters %v; expected %v", len(inst.Parameters), expected[i])
		}
	}
}

func TestGetInstruction(t *testing.T) {
	inputs := [][]int{
		[]int{1002, 0, 2, 0, 99},
		[]int{3, 225, 4, 225},
		[]int{4, 225, 1, 225}}
	expected := []Instruction{
		Instruction{
			Position: 0,
			Op:       MultiplyOp,
			Parameters: []Parameter{
				Parameter{Value: 0},
				Parameter{Value: 2},
				Parameter{Value: 0}},
			Modes: []ParameterMode{
				PositionMode,
				ImmediateMode,
				PositionMode}},
		Instruction{
			Position: 2,
			Op:       InputOp,
			Parameters: []Parameter{
				Parameter{Value: 225, Position: 3}},
			Modes: []ParameterMode{
				PositionMode}},
		Instruction{
			Position: 0,
			Op:       OutputOp,
			Parameters: []Parameter{
				Parameter{Value: 225}},
			Modes: []ParameterMode{
				PositionMode}}}

	expectedI := []int{4, 2, 2}

	for i, input := range inputs {
		inst := GetInstruction(0, input)

		if inst.NextInstruction() != expectedI[i] {
			t.Errorf("incorrect next position %v; expected %v", inst.NextInstruction(), expectedI[i])
		}

		if inst.Op != expected[i].Op {
			t.Errorf("incorrect op code %v; expected %v", inst.Op, expected[i].Op)
		}

		if len(inst.Parameters) != len(expected[i].Parameters) {
			t.Errorf("incorrect parameter length %v; expected %v", len(inst.Parameters), len(expected[i].Parameters))
		}

		for j, p := range inst.Parameters {
			if p.Value != expected[i].Parameters[j].Value {
				t.Errorf("incorrect parameter %v; expected %v", p.Value, expected[i].Parameters[j].Value)
			}

			if inst.Modes[j] != expected[i].Modes[j] {
				t.Errorf("incorrect mode %v; expected %v", inst.Modes[j], expected[i].Modes[j])
			}
		}
	}
}

func TestFullIntegration(t *testing.T) {
	inputs := [][]int{
		[]int{3, 0, 4, 0, 99},
		[]int{3, 0, 1, 0, 0, 0, 4, 0, 99},
		[]int{3, 0, 2, 0, 0, 0, 4, 0, 99},
		[]int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
		[]int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8},
		[]int{3, 3, 1108, -1, 8, 3, 4, 3, 99},
		[]int{3, 3, 1107, -1, 8, 3, 4, 3, 99},
		[]int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
		[]int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1},
		[]int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31, 1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104, 999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99}}
	expected := []string{"4", "8", "16", "0", "1", "0", "1", "1", "1", "999"}

	for i, set := range inputs {
		position := 0
		br := bytes.NewReader([]byte{'4', '\n'})
		reader := bufio.NewReader(br)
		writer := bytes.NewBuffer([]byte{})
		instruction := GetInstruction(position, set)
		position = ExecInstruction(reader, writer, instruction, set)
		if position >= 0 {
			Process(reader, writer, position, set)
		}

		result := strings.Split(strings.TrimRight(fmt.Sprintf("%v", writer), "\n"), "\n")
		if result[len(result)-1] != expected[i] {
			t.Errorf("incorrect result %v for test %v; expected %v", result[len(result)-1], i+1, expected[i])
		}
	}
}
