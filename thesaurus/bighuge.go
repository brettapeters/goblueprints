package thesaurus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// BigHuge type houses the necessary API key and provides the
// Synonyms method that accesses the endpoint and parses the response.
type BigHuge struct {
	APIKey string
}

type synonyms struct {
	Noun *words `json:"noun"`
	Verb *words `json:"verb"`
}

type words struct {
	Syn []string `json:"syn"`
}

// Synonyms accesses the Big Huge Thesaurus endpoint and parses the results.
func (b *BigHuge) Synonyms(word string) ([]string, error) {
	var syns []string
	url := fmt.Sprintf("http://words.bighugelabs.com/api/2/%s/%s/json", b.APIKey, word)
	response, err := http.Get(url)
	if err != nil {
		return syns, err
	}
	var data synonyms
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return syns, err
	}
	syns = append(syns, data.Noun.Syn...)
	syns = append(syns, data.Verb.Syn...)
	return syns, nil
}
