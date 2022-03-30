package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type seed struct {
	id   string
	pass bool
}

// lol
// /Users/edward.sansome/code/slosh/slosh /path/to/rspec/file 1
// slosh [cmd] [loop]
func main() {


	// do this properly in the morning pls
	// db := flag.Bool("db", true, "speed up specs if you don't need the db :)")
	// flag.Parse()
	s := os.Args[1:]
	if len(s) != 2 {
		fmt.Println("usage: slosh [path/to/rspec/file] [loop]")
		os.Exit(1)
	}

	path := s[0]
	loop, err := strconv.Atoi(s[1])
	if err != nil {
		fmt.Println("usage: slosh [path/to/rspec/file] [loop]")
		os.Exit(1)
	}

	c := make(chan seed, loop)

	// add a cool msg for the user here

	for i := 0; i < cap(c); i++ {
		// reduce the chance of mysql deadlock hehe
		// we plus three as i could be 0 lol
		// this needs work
		// check if we need the db doe
		// if i > 0  {
		// 	time.Sleep(time.Second * time.Duration(i+2))
		// }
		go runspec(path, c)
	}

	for i := 0; i < cap(c); i++ {
		s := <-c
		if s.pass {
			color.Green(s.id)
		} else {
			color.Red(s.id)
		}
	}
}

func runspec(path string, c chan seed) {
	cmd := exec.Command("rspec", path)
	o, err := cmd.CombinedOutput()

	if err != nil {
		c <- seed{id: string(extractSeed(string(o))), pass: false}
		return
	}
	c <- seed{id: string(extractSeed(string(o))), pass: true}
}

func extractSeed(s string) string {
	st := strings.Split(s, "Randomized with seed ")
	sa := strings.Split(st[1], "\n")
	return fmt.Sprintf("seed %s\n", sa[0])
}