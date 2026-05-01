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
	opEND   = byte(zeroValue)
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
	stack     []int     // loops address stack.
	messenger chan rune // communication channel.
	offset    int       // bytecode offset.
	err       error     // compiler error state.
}

// This function runs the compiler concurrently.
func (c *compiler) run() {
	go c.l.run()
	for generator := generators[instr]; generator != nil; generator = generator(c) {
	}
}

// This function does some checks after compiling the ] operator.
func (c *compiler) endLoop(op byte) generator {
	if op == opEND {
		c.bytecode = append(c.bytecode, op)
		c.offset++
		c.err = fmt.Errorf("found EOF but at least one [ has not its matching ]")
		return nil
	} else {
		c.bytecode = append(c.bytecode, op)
		c.bytecode = append(c.bytecode, zeroValue, zeroValue)
		c.offset += endLoopOffset
		address := c.stack[len(c.stack)-1]
		binary.BigEndian.PutUint16(c.bytecode[c.offset-1:], uint16(address))
		binary.BigEndian.PutUint16(c.bytecode[address+1:], uint16(c.offset+1))
		c.stack = c.stack[:len(c.stack)-1]
		return generators[instr]
	}
}

// creates a new compiler.
func newCompiler(code string) *compiler {
	c := &compiler{
		bytecode:  make([]byte, zeroValue),
		messenger: make(chan rune),
		offset:    offsetInit,
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
				c.bytecode = append(c.bytecode, zeroValue, zeroValue)
				c.offset += uint16ByteSize
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
			var op byte
			for op = byte(<-c.messenger); op != opRB && op != opEND; op = byte(<-c.messenger) {
				c.bytecode = append(c.bytecode, op)
				c.offset++
				if op == opLB {
					c.stack = append(c.stack, c.offset)
					c.bytecode = append(c.bytecode, zeroValue, zeroValue)
					c.offset += uint16ByteSize
					if result := generators[loop](c); result == nil {
						return nil
					}
				}
			}
			return c.endLoop(op)
		},
	}
}
