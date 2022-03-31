package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

type seed struct {
	id   string
	pass bool
}

func main() {
	db := flag.Bool("db", true, "speed up specs if you don't need the db :)")
	path := flag.String("path", "", "path to the rspec file you want to slosh")
	loop := flag.Int("loop", 3, "how many times you would like to run the specs")
	flag.Parse()

	if *path == "" {
		fmt.Println("please supply path to spec file")
		os.Exit(1)
	}

	c := make(chan seed, *loop)

	// add a cool msg for the user here

	for i := 0; i < cap(c); i++ {
		// reduce the chance of mysql deadlock hehe
		// we plus three as i could be 0 lol
		// this should be improved as super hacky
		if i > 0 && *db {
			time.Sleep(time.Second * time.Duration(i+2))
		}
		go runspec(*path, c)
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
