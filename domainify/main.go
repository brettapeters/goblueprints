package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"

func main() {
	// seed the random source
	rand.Seed(time.Now().UTC().UnixNano())

	// read the tlds from the command line flag
	tlds := flag.String("tlds", "com, net", "top level domains")
	flag.Parse()
	// strip whitespace
	domains := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, *tlds)
	// split on commas
	tldList := strings.Split(domains, ",")

	// read in the input from STDIN
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		// convert the input string to lowercase
		text := strings.ToLower(s.Text())
		var newText []rune
		// range over each character in the string and replace/skip
		// if necessary. Afterwards, stick a random tld to the end
		for _, r := range text {
			if unicode.IsSpace(r) {
				r = '-'
			}
			if !strings.ContainsRune(allowedChars, r) {
				continue
			}
			newText = append(newText, r)
		}
		fmt.Println(string(newText) + "." + tldList[rand.Intn(len(tldList))])
	}
}
