package dictionary

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Write dictionary to file.
func Write(path string, dictonary map[string]string) error {
	cards := []string{}
	for k, v := range dictonary {
		cards = append(cards, fmt.Sprintf("[%s|%s]", k, v))
	}
	data := strings.Join(cards, "\n")
	return ioutil.WriteFile(path, []byte(data), 0644)
}
