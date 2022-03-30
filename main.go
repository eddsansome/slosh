package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type seed struct {
	id string
	pass bool
}

// lol
// /Users/edward.sansome/code/slosh/slosh /path/to/rspec/file 1
// slosh [cmd] [loop]
func main() {

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

	c := make(chan seed)

	for i := 0; i < loop; i++ {
		// reduce the chance of mysql deadlock hehe
		// we plus three as i could be 0 lol
		// this needs work
		time.Sleep(time.Second * time.Duration(i+3))
		go runspec(path, c)
	}
	for {
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
	o, err := cmd.Output()
	if err != nil {
	c <- seed{id: string(extractSeed(err.Error())), pass: false} 
	}
	c <- seed{id: string(extractSeed(string(o))), pass: true} 
}

func extractSeed(s string) string {
	st := strings.Split(s, "Randomized with seed ")
	sa := strings.Split(st[1], "\n")
	return fmt.Sprintf("seed %s\n", sa[0])
}
