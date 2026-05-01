package brainfact

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

// Globalconfigurable variables.
var (
	TapeSize uint = math.MaxInt16
	Prompt        = "bf > "
)

// interpreter compiles and executes brainfuck code.
type interpreter struct {
	code     string
	compiler *compiler
	tape     []byte
	reader   *bufio.Reader
}

// Run creates a new interpreter and runs the code passed as argument. It returns nil if no error happened.
func Run(code string) error {
	i := &interpreter{
		code:     code,
		compiler: newCompiler(code),
		tape:     make([]byte, TapeSize),
		reader:   bufio.NewReader(os.Stdin),
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
	bcLimit := uint(len(bytecode))
	var pointer uint
	var pc uint
	var op byte
	for {
		if pc >= bcLimit {
			return fmt.Errorf("progam counter out of bounds: %d", pc)
		}
		op = bytecode[pc]
		switch op {
		case opLeft:
			if pointer == 0 {
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
			if pc+4 >= bcLimit {
				return fmt.Errorf("malformed bytecode: truncated jump address at %d", pc)
			}
			if tape[pointer] == 0 {
				pc = uint(binary.BigEndian.Uint32(bytecode[pc+1:]))
			} else {
				pc += 5
			}
		case opRB:
			if pc+4 >= bcLimit {
				return fmt.Errorf("malformed bytecode: truncated jump address at %d", pc)
			}
			if tape[pointer] != 0 {
				pc = uint(binary.BigEndian.Uint16(bytecode[pc+1:]))
			} else {
				pc += 5
			}
		case opComma:
			b, err := i.reader.ReadByte()
			if err == nil {
				tape[pointer] = b
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
