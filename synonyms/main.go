package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	if len(apiKey) == 0 {
		log.Fatalln("Missing API key for Big Huge Thesaurus")
	}
	thesaurus := thesaurus.BigHuge{APIKey: os.Getenv("BHT_APIKEY")}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalln("Failed when looking for synonyms for \""+word+"\"", err)
		}
		if len(syns) == 0 {
			log.Fatalln("Couldn't find any synonyms for \"" + word + "\"")
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
