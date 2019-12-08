package intcode

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// OpCode ...
type OpCode int

// ParameterMode ...
type ParameterMode int

const (
	// AddOp ...
	AddOp OpCode = 1
	// MultiplyOp ...
	MultiplyOp OpCode = 2
	// InputOp ...
	InputOp OpCode = 3
	// OutputOp ...
	OutputOp OpCode = 4
	// JumpTrue ...
	JumpTrue OpCode = 5
	// JumpFalse ...
	JumpFalse OpCode = 6
	// LessThan ...
	LessThan OpCode = 7
	// Equals ...
	Equals OpCode = 8
	// TerminateOp ...
	TerminateOp OpCode = 99
	// PositionMode ...
	PositionMode ParameterMode = 0
	// ImmediateMode ...
	ImmediateMode ParameterMode = 1
)

// Executable ...
type Executable interface {
	OpCode() OpCode
	GetValue(int, []int) int
	GetInput(io.Reader, io.Writer, string) (int, error)
	SetValue(int, int, []int) []int
	SetOutput(io.Writer, string)
	SetNext(int)
	NextInstruction() int
}

// Instruction ...
type Instruction struct {
	Position   int
	Op         OpCode
	Parameters []Parameter
	Modes      []ParameterMode
	Next       *int
}

// OpCode ...
func (i *Instruction) OpCode() OpCode {
	return i.Op
}

// GetValue ...
func (i *Instruction) GetValue(index int, set []int) int {
	mode := i.Modes[index]
	parm := i.Parameters[index]
	switch mode {
	case PositionMode:
		return set[parm.Value]
	case ImmediateMode:
		return parm.Value
	}

	return 0
}

// GetInput ...
func (i *Instruction) GetInput(in *bufio.Reader, out io.Writer, message string) (int, error) {
	text, _ := in.ReadString('\n')
	t := strings.TrimRightFunc(text, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	return strconv.Atoi(t)
}

// SetValue ...
func (i *Instruction) SetValue(value int, index int, set []int) []int {
	mode := i.Modes[index]
	parm := i.Parameters[index]
	switch mode {
	case PositionMode:
		set[parm.Value] = value
	case ImmediateMode:
		set[parm.Position] = value
	}

	return set
}

// SetOutput ...
func (i *Instruction) SetOutput(out io.Writer, message string) {
	fmt.Fprintf(out, "%s\n", message)
}

// SetNext ...
func (i *Instruction) SetNext(position int) {
	i.Next = &position
}

// NextInstruction ...
func (i *Instruction) NextInstruction() int {
	switch i.Op {
	case TerminateOp:
		return -1
	default:
		if i.Next != nil {
			return *i.Next
		}
		return i.Position + len(i.Parameters) + 1
	}
}

// Parameter ...
type Parameter struct {
	Value    int
	Position int
}

// SplitOp ...
func SplitOp(opcode int) (OpCode, []ParameterMode) {
	padded := []byte(fmt.Sprintf("%6v", strconv.Itoa(opcode)))

	g := strings.TrimLeftFunc(string(padded[4:]), func(r rune) bool {
		return unicode.IsSpace(r)
	})
	oc64, _ := strconv.ParseInt(g, 10, 64)
	op := int(oc64)

	params := []ParameterMode{}
	for i := 3; i >= 0; i-- {
		p, _ := strconv.Atoi(string(padded[i]))
		var mode ParameterMode
		switch p {
		case 1:
			mode = ImmediateMode
		case 0:
		default:
			mode = PositionMode
		}
		params = append(params, mode)
	}

	var oc OpCode
	switch op {
	case 1:
		oc = AddOp
	case 2:
		oc = MultiplyOp
	case 3:
		oc = InputOp
	case 4:
		oc = OutputOp
	case 5:
		oc = JumpTrue
	case 6:
		oc = JumpFalse
	case 7:
		oc = LessThan
	case 8:
		oc = Equals
	case 99:
		oc = TerminateOp
	default:
		panic(fmt.Sprintf("incorrect op code: %v; for %v", op, opcode))
	}

	return oc, params
}

// GetParameters ...
func GetParameters(quantity int, position int, set []int) []Parameter {
	params := make([]Parameter, quantity)

	for j, p := range set[position : position+len(params)] {
		params[j] = Parameter{Value: p, Position: j + position}
	}

	return params
}

// GetInstruction ...
func GetInstruction(i int, set []int) Instruction {
	code, modes := SplitOp(set[i])

	var params []Parameter
	var in Instruction

	switch code {
	case AddOp, MultiplyOp, LessThan, Equals:
		params = GetParameters(3, i+1, set)
	case JumpTrue, JumpFalse:
		params = GetParameters(2, i+1, set)
	case InputOp, OutputOp:
		params = GetParameters(1, i+1, set)
	case TerminateOp:
		params = []Parameter{}
	default:
		panic("unknown op code")
	}

	in = Instruction{
		Position:   i,
		Op:         code,
		Parameters: params,
		Modes:      modes}

	return in
}

// ExecInstruction ...
func ExecInstruction(in *bufio.Reader, out io.Writer, ex Instruction, set []int) int {
	switch ex.OpCode() {
	case AddOp:
		r := Add(ex.GetValue(0, set), ex.GetValue(1, set))
		set = ex.SetValue(r, 2, set)
	case MultiplyOp:
		r := Mult(ex.GetValue(0, set), ex.GetValue(1, set))
		set = ex.SetValue(r, 2, set)
	case InputOp:
		input, _ := ex.GetInput(in, out, "Enter input value: ")
		set = ex.SetValue(input, 0, set)
	case OutputOp:
		ex.SetOutput(out, fmt.Sprintf("%v", ex.GetValue(0, set)))
	case JumpTrue:
		if ex.GetValue(0, set) != 0 {
			ex.SetNext(ex.GetValue(1, set))
		}
	case JumpFalse:
		if ex.GetValue(0, set) == 0 {
			ex.SetNext(ex.GetValue(1, set))
		}
	case LessThan:
		a := ex.GetValue(0, set)
		b := ex.GetValue(1, set)
		if a < b {
			ex.SetValue(1, 2, set)
		} else {
			ex.SetValue(0, 2, set)
		}
	case Equals:
		a := ex.GetValue(0, set)
		b := ex.GetValue(1, set)
		if a == b {
			ex.SetValue(1, 2, set)
		} else {
			ex.SetValue(0, 2, set)
		}
	case TerminateOp:
		fallthrough
	default:
		return ex.NextInstruction()
	}

	return ex.NextInstruction()
}

// Add ...
func Add(i, j int) int {
	return i + j
}

// Mult ...
func Mult(i, j int) int {
	return i * j
}

// ReadCodes ...
func ReadCodes(path string) []int {
	codes := []int{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), ",")
		for _, v := range values {
			vi, _ := strconv.Atoi(v)
			codes = append(codes, vi)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return codes
}

// Process ...
func Process(in *bufio.Reader, out io.Writer, position int, set []int) {
	instruction := GetInstruction(position, set)
	position = ExecInstruction(in, out, instruction, set)
	if position >= 0 {
		Process(in, out, position, set)
	}
}
