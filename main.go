package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	inputPath := flag.String("i", "", "input file")
	dump := flag.Bool("dump", false, "dump parsed instructions")
	repl := flag.Bool("repl", false, "run interactive REPL")
	flag.Parse()

	if *repl {
		startRepl()
		return
	}

	if len(*inputPath) == 0 {
		log.Fatalln("ERROR: Not input file")
		return
	}

	file, err := os.Open(*inputPath)
	if err != nil {
		log.Fatalln("ERROR:", err)
		return
	}
	defer file.Close()

	program := Program{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lineTokens := strings.Split(line, " ")
		if err := program.Parse(lineTokens); err != nil {
			log.Fatalln("ERROR:", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln("ERROR:", err)
		return
	}

	if *dump {
		program.dump()
	}

	program.run()
}

func startRepl() {
	fmt.Println("Interactive REPL. Type `exit` to quit.")
	program := Program{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" {
			break
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		tokens := strings.Split(line, " ")
		err := program.Parse(tokens)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}
		stack := program.run()
		if len(stack) > 0 {
			fmt.Printf(":- ")
			for i, n := range stack {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%v", n)
			}
			fmt.Printf("\n")
		}
	}
}

type Program []ParsedLine

func (program Program) dump() {
	for _, pl := range program {
		fmt.Printf("Line %d:\n", pl.line)
		for _, ins := range pl.instrs {
			fmt.Printf("    %s\n", ins.String())
		}
	}
}

func (program Program) run() []int64 {
	stack := []int64{}
	for _, pl := range program {
		for _, ins := range pl.instrs {
			switch ins.kind {
			case PUSH_INT:
				stack = append(stack, ins.operand)
			case PLUS:
				if len(stack) < 2 {
					log.Fatalf("ERROR: not enough values on stack for '%s'\n", ins.kind.String())
				}
				b := stack[len(stack)-1]
				a := stack[len(stack)-2]
				stack = stack[:len(stack)-2]
				stack = append(stack, a+b)
			case PRINT:
				if len(stack) < 1 {
					log.Fatalln("ERROR: empty stack on 'print'")
				}
				a := stack[len(stack)-1]
				stack = stack[0 : len(stack)-1]
				println(a)
			default:
				log.Fatalln("ERROR: unimplemented", ins.kind.String())
				return stack
			}
		}
	}
	return stack
}

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
		case "-":
			{
				pl.instrs = append(pl.instrs, Instr{MINUS, 0})
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

func (ins Instr) String() string {
	if ins.kind.HasOperand() {
		return fmt.Sprintf("%s(%v)", ins.kind.String(), ins.operand)
	} else {
		return fmt.Sprintf("%s", ins.kind.String())
	}
}

type InstrKind uint

const (
	PUSH_INT InstrKind = iota
	PLUS
	MINUS
	PRINT
)

func (ins InstrKind) CheckArity(stack_size int) bool {
	n, _ := ins.Arity()
	if stack_size < n {
		log.Fatalf("ERROR: not enough values on stack for '%s'\n", ins.String())
	}
}

func (ins InstrKind) Arity() (int, int) {
	switch ins {
	case PLUS, MINUS:
		return 2, 1
	case PRINT:
		return 1, 0
	case PUSH_INT:
		return 0, 1
	default:
		panic("Unknown InstrKind")
	}
}

func (ins InstrKind) String() string {
	switch ins {
	case PUSH_INT:
		return "push_int"
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case PRINT:
		return "print"
	default:
		panic("Unknown InstrKind")
	}
}

func (ins InstrKind) HasOperand() bool {
	switch ins {
	case PUSH_INT:
		return true
	default:
		return false
	}
}
