package expr

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
	"text/scanner"
)

type Instruction uint8

const (
	InstInvalid Instruction = iota
	InstFloatLiteral
	InstPlus
	InstMinus
	InstLog
)

func Run(insts []byte) float64 {
	var stack = make([]float64, 0, 16)

	pc := 0
	for pc < len(insts) {
		var inst Instruction
		inst, pc = Instruction(insts[pc]), pc+1
		switch inst {
		case InstFloatLiteral:
			val := math.Float64frombits(binary.LittleEndian.Uint64(insts[pc : pc+8]))
			pc += 8
			stack = append(stack, val)
		case InstPlus:
			stack[len(stack)-2] += stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case InstMinus:
			stack[len(stack)-2] -= stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case InstLog:
			stack[len(stack)-1] = math.Log(stack[len(stack)-1])
		default:
			log.Panic(inst)
		}
	}
	if len(stack) != 1 {
		panic("stack")
	}
	return stack[0]
}

type opSpec struct {
	nArg int
	inst Instruction
}

type Compiler struct {
	sc         *scanner.Scanner
	stackDepth int
	inst       []byte
	ops        map[string]opSpec
}

func (c *Compiler) addInst(i Instruction) {
	c.inst = append(c.inst, uint8(i))
}

func (c *Compiler) addFloatLiteral(val float64) {
	var buf [8]uint8
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(val))
	c.inst = append(c.inst, buf[:]...)
}

func (c *Compiler) compile() int {
	for {
		ch := c.sc.Scan()
		if ch == scanner.EOF {
			return -1
		}
		if ch == ')' {
			return ')'
		}
		if ch == '(' {
			if endch := c.compile(); endch != ')' {
				panic("))))")
			}
		}
		if ch == scanner.Int {
			val, err := strconv.ParseInt(c.sc.TokenText(), 0, 64)
			if err != nil {
				log.Panic(err)
			}
			c.addFloatLiteral(float64(val))
			c.stackDepth++
			continue
		}
		if ch == scanner.Float {
			val, err := strconv.ParseFloat(c.sc.TokenText(), 64)
			if err != nil {
				log.Panic(err)
			}
			c.addFloatLiteral(val)
			c.stackDepth++
			continue
		}
		var opName string
		if ch >= 1 && ch < 128 {
			opName = fmt.Sprintf("%c", ch)
		} else if ch == scanner.Ident {
			opName = c.sc.TokenText()
		}
		op, ok := c.ops[opName]
		if !ok {
			log.Panic(opName)
		}

		c.addInst(op.inst)
		if c.stackDepth < op.nArg {
			panic("narg")
		}
		c.stackDepth -= op.nArg
		c.stackDepth++
	}
}

func Compile(in io.Reader) []byte {
	c := Compiler{
		sc: &scanner.Scanner{},
		ops: map[string]opSpec{
			"+":   opSpec{nArg: 2, inst: InstPlus},
			"-":   opSpec{nArg: 2, inst: InstMinus},
			"log": opSpec{nArg: 1, inst: InstLog},
		},
	}

	c.sc.Error = func(_ *scanner.Scanner, msg string) {
		log.Panic(msg)
	}
	c.sc.Mode = scanner.GoTokens
	c.sc.Init(in)
	if endch := c.compile(); endch != -1 {
		panic("eof")
	}
	return c.inst
}
