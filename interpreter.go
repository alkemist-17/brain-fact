package brainfact

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

// Global Consts
const (
	jump           = 3
	firstCellIndex = 0
	zeroValue      = 0
	offsetInit     = -1
	uint16ByteSize = 2
	endLoopOffset  = 3
)

// Globalconfigurable variables.
var (
	TapeSize    uint = math.MaxInt16
	Prompt           = " > "
	InputPrompt      = " · "
)

// interpreter compiles and executes brainfuck code.
type interpreter struct {
	code     string
	compiler *compiler
	tape     []byte
	scanner  *bufio.Scanner
}

// Run creates a new interpreter and runs the code passed as argument. It returns nil if no error happened.
func Run(code string) error {
	i := &interpreter{
		code:     code,
		compiler: newCompiler(code),
		tape:     make([]byte, TapeSize),
		scanner:  bufio.NewScanner(os.Stdin),
	}
	i.compiler.run()
	if i.compiler.err != nil {
		return i.compiler.err
	}
	return i.run()
}

// run runs the brainfuck interpreter.
func (i *interpreter) run() error {
	tape := i.tape
	bytecode := i.compiler.bytecode
	tapeLimit := uint(len(tape))
	var pointer uint
	var pc uint
	var op byte
	for {
		op = bytecode[pc]
		switch op {
		case opLeft:
			if pointer == firstCellIndex {
				return fmt.Errorf("tape underflow")
			}
			pointer--
			pc++
		case opRight:
			if pointer == tapeLimit-1 {
				return fmt.Errorf("tape overflow")
			}
			pointer++
			pc++
		case opAdd:
			tape[pointer]++
			pc++
		case opSub:
			tape[pointer]--
			pc++
		case opLB:
			if tape[pointer] == zeroValue {
				pc = uint(binary.BigEndian.Uint16(bytecode[pc+1:]))
			} else {
				pc += jump
			}
		case opRB:
			if tape[pointer] != zeroValue {
				pc = uint(binary.BigEndian.Uint16(bytecode[pc+1:]))
			} else {
				pc += jump
			}
		case opComma:
			fmt.Printf("\n%s", InputPrompt)
			for i.scanner.Scan() {
				if b := i.scanner.Bytes(); len(b) > zeroValue {
					tape[pointer] = b[zeroValue]
				}
				break
			}
			pc++
		case opDot:
			fmt.Print(string(tape[pointer]))
			pc++
		case opEND:
			return nil
		default:
			return fmt.Errorf("unknown bytecode %v", op)
		}
	}
}
