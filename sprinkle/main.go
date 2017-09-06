package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

var transformsFile = "transforms.txt"

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// open the transforms txt file
	var (
		f   *os.File
		err error
		ex  string
	)
	if ex, err = os.Executable(); err != nil {
		fmt.Println(err)
		return
	}
	if f, err = os.Open(path.Join(path.Dir(ex), transformsFile)); err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// scan over each line and add it to the slice of strings
	// called 'transforms'
	var transforms []string
	ts := bufio.NewScanner(f)
	for ts.Scan() {
		transforms = append(transforms, ts.Text())
	}

	// scan over the input from STDIN and replace the * symbol
	// with the input word.
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, "*", s.Text(), -1))
	}
}
