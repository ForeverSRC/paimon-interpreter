package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/ForeverSRC/paimon-interpreter/pkg/repl"
)

const Welcome = "Hello %s! This is the Paimon programming language!\nFeel free to type in commands\n"

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Println(PAIMON)
	fmt.Printf(Welcome, usr)
	repl.Start(os.Stdin, os.Stdout)

}
