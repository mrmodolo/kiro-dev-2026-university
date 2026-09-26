package main

import (
	"fmt"
	"os"
	"strings"

	"helloworld-pbt/greeter"
)

// main prints a greeting. With no arguments it prints "Hello, World!".
// With arguments it greets the joined name, e.g. `helloworld-pbt Ada Lovelace`.
func main() {
	name := strings.Join(os.Args[1:], " ")
	fmt.Println(greeter.Greet(name))
}
