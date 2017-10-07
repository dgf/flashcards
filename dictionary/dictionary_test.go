package dictionary_test

import (
	"io/ioutil"
	"testing"

	"github.com/dgf/flashcards/dictionary"
)

func TestDictonaryParser(t *testing.T) {
	dict := map[string]string{
		"uno":  "eins",
		"dos":  "zwei",
		"tres": "drei",
		`1\n2`: `a\nb`,
	}

	// create test file
	file, err := ioutil.TempFile("", "test.dict")
	if err != nil {
		t.Error(err)
	}
	path := file.Name()

	// close test file
	if err := file.Close(); err != nil {
		t.Error(err)
	}

	// write dictonary file
	if err := dictionary.Write(path, dict); err != nil {
		t.Error(err)
	}
	t.Log(path)

	// read dictionary file
	act, err := dictionary.Read(path)
	if err != nil {
		t.Error(err)
	}

	// assert
	for front, back := range dict {
		if back != act[front] {
			t.Errorf("invalid write/read cycle\n\tEXP: [%s|%s]\n\tACT: [%s|%s]", front, back, front, act[front])
		}
	}
}
