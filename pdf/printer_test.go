package pdf_test

import (
	"io/ioutil"
	"testing"

	"github.com/dgf/flashcards/config"
	"github.com/dgf/flashcards/pdf"
	"github.com/dgf/flashcards/testdata"
)

func TestPrint(t *testing.T) {

	// create test file
	file, err := ioutil.TempFile("", "flashcards.pdf")
	if err != nil {
		t.Fatal(err)
	}
	name := file.Name()
	_ = file.Close()

	// create PDF and print cards
	if err := pdf.New(name, pdf.Options{
		Size: "A6",
		Card: config.Card{
			Unit:   "mm",
			Width:  73,
			Height: 31,
			Margin: 2,
		},
		Font: config.Font{
			Size: 37,
			Name: "Arial",
		},
	}).Print(testdata.Cards); err != nil {
		t.Fatal(err)
	}

	// read test file
	b, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	// assert test file format
	exp := "%PDF-1.3"
	act := string(b[0:8])
	if exp != act {
		t.Errorf("invalid file format\n\tEXP: %s\n\tACT: %s", exp, act)
	}
}
