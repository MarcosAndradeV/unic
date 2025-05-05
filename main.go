package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	tempFile := "temp/main.un"
	file, err := os.Open(tempFile)
	if err != nil {
		log.Fatalln("ERROR:", err)
		return
	}
	defer file.Close()

	program := Program{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineTokens := strings.Split(line, " ")
		slices.Reverse(lineTokens)
		if err := program.Parse(lineTokens); err != nil {
			log.Fatalln("ERROR:", err)
			return
		}
	}

	for _, pl := range program {
		for _, ins := range pl.instrs {
			fmt.Printf("ins: %s\n", ins.String())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("ERROR:", err)
		return
	}
}

type Program []ParsedLine

func (p *Program) Parse(lineTokens []string) error {
	pl := ParsedLine{line: len(*p), instrs: []Instr{}}
	for _, token := range lineTokens {
		switch token {
		case "print":
			{
				pl.instrs = append(pl.instrs, Instr{PRINT, 0})
			}
		case "+":
			{
				pl.instrs = append(pl.instrs, Instr{PLUS, 0})
			}
		default:
			{
				if i, err := strconv.ParseInt(token, 10, 64); err == nil {
					pl.instrs = append(pl.instrs, Instr{PUSH_INT, i})
					continue
				}

				return fmt.Errorf("Not implemented token: `%s`", token)
			}
		}
	}
	*p = append(*p, pl)
	return nil
}

type ParsedLine struct {
	line   int
	instrs []Instr
}

type Instr struct {
	kind    InstrKind
	operand int64
}

func (ins *Instr) String() string {
	if ins.kind.HasOperand() {
		return fmt.Sprintf("%s(%v)", ins.kind.String(), ins.operand)
	} else {
		return fmt.Sprintf("%s", ins.kind.String())
	}
}

type InstrKind uint

func (ins *InstrKind) String() string {
	switch *ins {
	case PUSH_INT:
		return "push_int"
	case PLUS:
		return "+"
	case PRINT:
		return "print"
	default:
		return "Unknow"
	}
}

func (ins *InstrKind) HasOperand() bool {
	switch *ins {
	case PUSH_INT:
		return true
	case PLUS:
		return false
	case PRINT:
		return false
	default:
		return false
	}
}

const (
	PUSH_INT InstrKind = iota
	PLUS
	PRINT
)
