package brainfact

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// ISA
const (
	opLeft  = byte('<')
	opRight = byte('>')
	opAdd   = byte('+')
	opSub   = byte('-')
	opComma = byte(',')
	opDot   = byte('.')
	opLB    = byte('[')
	opRB    = byte(']')
	opEND   = byte(0)
)

// Lexer
type lexer struct {
	reader    *strings.Reader // strings reader.
	messenger chan rune       // communication channel.
}

// creates a new lexer.
func newLexer(code string, messenger chan rune) *lexer {
	return &lexer{
		reader:    strings.NewReader(code),
		messenger: messenger,
	}
}

// This function scanns the input code concurrently.
func (l *lexer) run() {
	for r, _, e := l.reader.ReadRune(); e == nil; r, _, e = l.reader.ReadRune() {
		switch r {
		case rune(opLeft),
			rune(opRight),
			rune(opAdd),
			rune(opSub),
			rune(opComma),
			rune(opLB),
			rune(opRB),
			rune(opDot):
			l.messenger <- r
		}
	}
	l.messenger <- rune(opEND)
	close(l.messenger)
}

// Compiler
type compiler struct {
	l         *lexer    // lexical analyzer.
	bytecode  []byte    // compiled bytecode.
	stack     []int64   // loops address stack.
	messenger chan rune // communication channel.
	offset    int64     // bytecode offset.
	err       error     // error state of the compiler.
}

// This function runs the compiler concurrently.
func (c *compiler) run() {
	go c.l.run()
	for generator := generators[instr]; generator != nil; generator = generator(c) {
	}
}

// creates a new compiler.
func newCompiler(code string) *compiler {
	c := &compiler{
		bytecode:  make([]byte, 0),
		messenger: make(chan rune),
		offset:    -1,
	}
	c.l = newLexer(code, c.messenger)
	return c
}

// A generator is a function that maps runes to bytecode.
type generator func(*compiler) generator

// Types of generators.
const (
	instr = iota
	loop
)

// The generators slice.
var generators []generator

// Initializes the code generators functions.
func init() {
	generators = []generator{
		instr: func(c *compiler) generator {
			op := byte(<-c.messenger)
			c.bytecode = append(c.bytecode, op)
			c.offset++
			switch op {
			case opLB:
				c.stack = append(c.stack, c.offset)
				c.bytecode = append(c.bytecode, 0)
				c.bytecode = append(c.bytecode, 0)
				c.offset += 2
				return generators[loop]
			case opRB:
				c.bytecode = append(c.bytecode, opEND)
				c.offset++
				c.err = fmt.Errorf("found ] with not matching [")
				return nil
			case opEND:
				return nil
			default:
				return generators[instr]
			}
		},
		loop: func(c *compiler) generator {
			for {
				op := byte(<-c.messenger)
				switch op {
				case opEND:
					c.bytecode = append(c.bytecode, op)
					c.offset++
					c.err = fmt.Errorf("found EOF but at least one [ has not its matching ]")
					return nil
				case opRB:
					if len(c.stack) == 0 {
						c.bytecode = append(c.bytecode, opEND)
						c.offset++
						c.err = fmt.Errorf("found ] with no matching [")
						return nil
					}
					c.bytecode = append(c.bytecode, op)
					c.bytecode = append(c.bytecode, 0, 0, 0, 0)
					c.offset += 5
					address := c.stack[len(c.stack)-1]
					c.stack = c.stack[:len(c.stack)-1]
					binary.BigEndian.PutUint32(c.bytecode[c.offset-3:], uint32(address))
					binary.BigEndian.PutUint32(c.bytecode[address+1:], uint32(c.offset+1))
					if len(c.stack) == 0 {
						return generators[instr]
					}
				case opLB:
					c.bytecode = append(c.bytecode, op)
					c.offset++
					c.stack = append(c.stack, c.offset)
					c.bytecode = append(c.bytecode, 0, 0, 0, 0)
					c.offset += 4
				default:
					c.bytecode = append(c.bytecode, op)
					c.offset++
				}
			}
		},
	}
}
