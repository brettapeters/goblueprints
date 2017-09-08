package thesaurus

type Thesaurus interface {
	Synonyms(word string) ([]string, error)
}
