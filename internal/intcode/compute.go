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
	// RelativeBase ...
	RelativeBase OpCode = 9
	// TerminateOp ...
	TerminateOp OpCode = 99
	// PositionMode ...
	PositionMode ParameterMode = 0
	// ImmediateMode ...
	ImmediateMode ParameterMode = 1
	// RelativeMode ...
	RelativeMode ParameterMode = 2
)

// Instruction ...
type Instruction struct {
	Position   int
	RelPos     int
	Op         OpCode
	Parameters []Parameter
	Modes      []ParameterMode
	Next       *int
	DataSet    []int
	Input      *bufio.Reader
	Output     io.Writer
}

func (i *Instruction) getValue(index int) (int, error) {
	mode := i.Modes[index]
	parm := i.Parameters[index]
	switch mode {
	case PositionMode:
		exp := parm.Value + 1
		if len(i.DataSet) < exp {
			n := make([]int, exp)
			copy(n, i.DataSet)
			i.DataSet = n
		}
		return i.DataSet[parm.Value], nil
	case ImmediateMode:
		return parm.Value, nil
	case RelativeMode:
		exp := i.RelPos + parm.Value + 1
		if len(i.DataSet) < exp {
			n := make([]int, exp)
			copy(n, i.DataSet)
			i.DataSet = n
		}
		return i.DataSet[i.RelPos+parm.Value], nil
	}

	return 0, &CodeTerminationError{exitCode: 1, message: "unknown parameter mode"}
}

func (i *Instruction) getInput() (int, error) {
	text, err := i.Input.ReadString('\n')
	if err != nil {
		return 0, err
	}

	t := strings.TrimRightFunc(text, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	return strconv.Atoi(t)
}

func (i *Instruction) setValue(index int, value int) error {
	mode := i.Modes[index]
	parm := i.Parameters[index]

	switch mode {
	case PositionMode:
		exp := parm.Value + 1
		if len(i.DataSet) < exp {
			n := make([]int, exp)
			copy(n, i.DataSet)
			i.DataSet = n
		}
		i.DataSet[parm.Value] = value
		return nil
	case ImmediateMode:
		exp := parm.Position + 1
		if len(i.DataSet) < exp {
			n := make([]int, exp)
			copy(n, i.DataSet)
			i.DataSet = n
		}
		i.DataSet[parm.Position] = value
		return nil
	case RelativeMode:
		exp := i.RelPos + parm.Value + 1
		if len(i.DataSet) < exp {
			n := make([]int, exp)
			copy(n, i.DataSet)
			i.DataSet = n
		}
		i.DataSet[i.RelPos+parm.Value] = value
		return nil
	}

	return &CodeTerminationError{exitCode: 1, message: "unknown parameter mode"}
}

func (i *Instruction) setOutput(out int) {
	fmt.Fprintf(i.Output, "%v\n", out)
}

func (i *Instruction) setNext(position int) {
	i.Next = &position
}

func (i *Instruction) nextPosition() (int, error) {
	switch i.Op {
	case TerminateOp:
		return 0, &CodeTerminationError{exitCode: 0, message: "success"}
	default:
		if i.Next != nil {
			return *i.Next, nil
		}
		return i.Position + len(i.Parameters) + 1, nil
	}
}

func (i *Instruction) loadParams(code OpCode) (*[]Parameter, error) {
	var params []Parameter

	switch code {
	case AddOp, MultiplyOp, LessThan, Equals:
		params = GetParameters(3, i.Position+1, i.DataSet)
	case JumpTrue, JumpFalse:
		params = GetParameters(2, i.Position+1, i.DataSet)
	case InputOp, OutputOp, RelativeBase:
		params = GetParameters(1, i.Position+1, i.DataSet)
	case TerminateOp:
		params = []Parameter{}
	default:
		return nil, &CodeTerminationError{exitCode: 1, message: "unknown opcode"}
	}

	return &params, nil
}

// Step ...
func (i *Instruction) Step() error {
	pos, err := i.nextPosition()
	if err != nil {
		return err
	}

	code, modes := SplitOp(i.DataSet[pos])
	i.Op = code
	i.Position = pos
	i.Next = nil
	p, err := i.loadParams(i.Op)
	if err != nil {
		return err
	}

	i.Parameters = *p
	i.Modes = modes
	if err != nil {
		return nil
	}

	return nil
}

// Exec ...
func (i *Instruction) Exec() error {
	switch i.Op {
	case AddOp:
		v1, err := i.getValue(0)
		if err != nil {
			return err
		}

		v2, err := i.getValue(1)
		if err != nil {
			return err
		}

		r := Add(v1, v2)
		if err != nil {
			return err
		}

		err = i.setValue(2, r)
		if err != nil {
			return err
		}
	case MultiplyOp:
		v1, err := i.getValue(0)
		if err != nil {
			return err
		}

		v2, err := i.getValue(1)
		if err != nil {
			return err
		}

		r := Mult(v1, v2)

		err = i.setValue(2, r)
		if err != nil {
			return err
		}
	case InputOp:
		input, err := i.getInput()
		if err != nil {
			return err
		}

		i.setValue(0, input)
	case OutputOp:
		v1, err := i.getValue(0)
		if err != nil {
			return err
		}

		i.setOutput(v1)
	case JumpTrue:
		v1, err := i.getValue(0)
		if err != nil {
			return err
		}

		// jump to position if value is not zero
		if v1 != 0 {
			v2, err := i.getValue(1)
			if err != nil {
				return err
			}

			i.setNext(v2)
		}
	case JumpFalse:
		v1, err := i.getValue(0)
		if err != nil {
			return err
		}

		// jump to position if value is zero
		if v1 == 0 {
			v2, err := i.getValue(1)
			if err != nil {
				return err
			}

			i.setNext(v2)
		}
	case LessThan:
		v1, err := i.getValue(0)
		if err != nil {
			return err
		}

		v2, err := i.getValue(1)
		if err != nil {
			return err
		}

		// if first value is less than second value, set third value to 1 else 0
		if v1 < v2 {
			i.setValue(2, 1)
		} else {
			i.setValue(2, 0)
		}
	case Equals:
		v1, err := i.getValue(0)
		if err != nil {
			return err
		}

		v2, err := i.getValue(1)
		if err != nil {
			return err
		}

		// if first value is equal to second value, set third value to 1 else 0
		if v1 == v2 {
			i.setValue(2, 1)
		} else {
			i.setValue(2, 0)
		}
	case RelativeBase:
		v1, err := i.getValue(0)
		if err != nil {
			return err
		}

		i.RelPos = i.RelPos + v1
	case TerminateOp:
		return &CodeTerminationError{exitCode: 0, message: "success"}
	default:
		m := fmt.Sprintf("unknown instruction code %v", i.Op)
		return &CodeTerminationError{exitCode: 1, message: m}
	}

	return nil
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
		case 2:
			mode = RelativeMode
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
	case 9:
		oc = RelativeBase
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

// CodeTerminationError ...
type CodeTerminationError struct {
	exitCode int
	message  string
}

// Error ...
func (e *CodeTerminationError) Error() string {
	return fmt.Sprintf("execution terminated with exit code %v; %s", e.exitCode, e.message)
}

func newInstructionSet(startPosition int, dataSet []int) *Instruction {
	inst := &Instruction{
		Position: -1,
		DataSet:  dataSet,
		Next:     &startPosition}

	return inst
}

// Process ...
func Process(in io.Reader, out io.Writer, position int, set []int) {
	comp := newInstructionSet(position, set)
	reader := bufio.NewReader(in)

	comp.Input = reader
	comp.Output = out

	var err error
	for {
		err = comp.Step()
		if err != nil {
			break
		}

		err = comp.Exec()
		if err != nil {
			break
		}
	}
}
