package dictionary

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var cardDelim = regexp.MustCompile(`\][ \n\r\t\f]*\[`)

// Read file and split dictionary cards [front|back]...[dos|zwei]
func Read(path string) (map[string]string, error) {
	dict := map[string]string{}

	// read file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := string(b)

	// trim first and last card markers
	data = strings.Trim(data, "[ ]\n")

	// split cards with abritrary line breaks
	cards := cardDelim.Split(data, -1)

	// split and put front and back
	fails := []string{}
	for _, c := range cards {
		sides := strings.SplitN(c, "|", 2)
		if len(sides) != 2 {
			fails = append(fails, fmt.Sprintf("[%s]", c))
		}
		front := strings.Trim(sides[0], " \n\r\t\f")
		back := strings.Trim(sides[1], " \n\r\t\f")
		dict[front] = back
	}

	// check invalid cards
	if len(fails) > 0 {
		return dict, fmt.Errorf("invalid cards:\n %s", strings.Join(fails, "\n"))
	}

	return dict, nil
}
