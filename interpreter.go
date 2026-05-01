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
	firstCellIndex  = 0
	zero            = 0
	U16Bytes        = 2
	U32Bytes        = 4
	FileExt         = ".bf"
	CompiledFileExt = ".bfc"
)

// Globalconfigurable variables.
var (
	TapeSize    uint = math.MaxInt16
	Prompt           = " > "
	InputPrompt      = " · "
	AddrSize         = U16Bytes
)

// interpreter compiles and executes brainfuck code.
type interpreter struct {
	compiler *compiler
	tape     []byte
	scanner  *bufio.Scanner
}

// Run creates a new interpreter and runs the code passed as argument. It returns nil if no error happened.
func Run(code string) error {
	i := &interpreter{
		compiler: newCompiler(code),
		tape:     make([]byte, TapeSize),
		scanner:  bufio.NewScanner(os.Stdin),
	}
	i.compiler.compile()
	if i.compiler.err != nil {
		return i.compiler.err
	}
	return i.run(nil)
}

// Run compiled brainfuck bytecode
func RunBytecode(bytecode []byte) error {
	i := &interpreter{
		tape:    make([]byte, TapeSize),
		scanner: bufio.NewScanner(os.Stdin),
	}
	return i.run(bytecode)
}

func Compile(fileName string, code string) error {
	i := &interpreter{
		compiler: newCompiler(code),
	}
	i.compiler.compile()
	if i.compiler.err != nil {
		return i.compiler.err
	}
	file, err := os.Create(fileName + CompiledFileExt)
	if err != nil {
		return err
	}
	defer file.Close()
	err = binary.Write(file, binary.BigEndian, i.compiler.bytecode)
	if err != nil {
		return err
	}
	return nil
}

// run runs the brainfuck interpreter.
func (i *interpreter) run(bc []byte) error {
	tape := i.tape
	var bytecode []byte
	if bc != nil {
		bytecode = bc
	} else {
		bytecode = i.compiler.bytecode
	}
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
			if tape[pointer] == zero {
				if AddrSize == U16Bytes {
					pc = uint(binary.BigEndian.Uint16(bytecode[pc+1:]))
				} else {
					pc = uint(binary.BigEndian.Uint32(bytecode[pc+1:]))
				}
			} else {
				pc += uint(AddrSize) + 1
			}
		case opRB:
			if tape[pointer] != zero {
				if AddrSize == U16Bytes {
					pc = uint(binary.BigEndian.Uint16(bytecode[pc+1:]))
				} else {
					pc = uint(binary.BigEndian.Uint32(bytecode[pc+1:]))
				}
			} else {
				pc += uint(AddrSize) + 1
			}
		case opComma:
			fmt.Printf("\n%s", InputPrompt)
			if i.scanner.Scan() {
				if b := i.scanner.Bytes(); len(b) > zero {
					tape[pointer] = b[zero]
				}
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
