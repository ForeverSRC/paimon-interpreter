package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/ForeverSRC/paimon-interpreter/cmd/consts"
	"github.com/ForeverSRC/paimon-interpreter/pkg/evaluator"
	"github.com/ForeverSRC/paimon-interpreter/pkg/lexer"
	"github.com/ForeverSRC/paimon-interpreter/pkg/object"
	"github.com/ForeverSRC/paimon-interpreter/pkg/parser"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Println(consts.PAIMON)
	fmt.Printf(consts.Welcome, usr.Name)
	Start(os.Stdin, os.Stdout)
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, consts.PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "parser errors: \n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
