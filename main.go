package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type seed struct {
	id   string
	pass bool
}

func main() {
	db := flag.Bool("db", true, "speed up specs if you don't need the db :)")
	path := flag.String("path", "", "path to the rspec file you want to slosh")
	// we should probably set a max value for now, until we can sort out the deadlocking issue
	// otherwise each run will take n+3 seconds to complete, which will get super slow
	loop := flag.Int("loop", 3, "how many times you would like to run the specs")
	flag.Parse()

	if *path == "" {
		fmt.Println("please supply path to spec file")
		os.Exit(1)
	}

	c := make(chan seed, *loop)
	s := spinner.New(spinner.CharSets[14], 50*time.Millisecond)
	s.Start()

	fmt.Println("beep boop, starting up rspec")

	for i := 0; i < cap(c); i++ {
		// reduce the chance of mysql deadlock hehe
		// we plus three as i could be 0 lol
		// this should be improved as super hacky
		if i > 0 && *db {
			time.Sleep(time.Second * time.Duration(i+3))
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
	s.Stop()
}

func runspec(path string, c chan seed) {
	cmd := exec.Command("rspec", path)
	o, err := cmd.CombinedOutput()

	if err != nil {
		c <- seed{id: extractSeed(o), pass: false}
		return
	}
	c <- seed{id: extractSeed(o), pass: true}
}

func extractSeed(o []byte) string {
	s := string(o)
	st := strings.Split(s, "Randomized with seed ")
	sa := strings.Split(st[1], "\n")
	return fmt.Sprintf("seed %s\n", sa[0])
}
