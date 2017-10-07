package pdf

import (
	"log"
	"strings"

	"github.com/dgf/flashcards/api"
	"github.com/jung-kurt/gofpdf"
)

// PDF implements api.Printer
type PDF struct {
	Options
	File string
}

// New creates a new PDF Printer
func New(path string, o Options) api.Printer {
	return &PDF{
		File:    path,
		Options: o,
	}
}

// font specific text splitting by line width
func cardLines(fpdf *gofpdf.Fpdf, width float64, text string) (lines []string) {
	for _, l := range fpdf.SplitLines([]byte(text), width) {
		lines = append(lines, string(l))
	}
	return
}

// multiple string words counting (simple split by space)
func wordCount(lines ...string) int {
	c := 0
	for _, l := range lines {
		if len(l) > 0 {
			c++
		}
		c += strings.Count(l, " ")
	}
	return c
}

// decrease font until it matches the card size
func decreaseFont(fpdf *gofpdf.Fpdf, w, h float64, text string) (lineHeight, lineCount float64) {
	_, lineHeight = fpdf.GetFontSize()
	words := wordCount(strings.Replace(text, "\n", " ", -1))
	lines := cardLines(fpdf, w, text)
	lineCount = float64(len(lines))
	linesWords := wordCount(lines...)

	for (linesWords > words) || (lineHeight*lineCount) > h {
		fontSize, _ := fpdf.GetFontSize()
		fpdf.SetFontSize(fontSize - 1)
		_, lineHeight = fpdf.GetFontSize()
		lines = cardLines(fpdf, w, text)
		lineCount = float64(len(lines))
		linesWords = wordCount(lines...)
	}
	return
}

// print card with adjusted font size
func printCard(fpdf *gofpdf.Fpdf, x, y, w, h, m float64, text string) {
	// replace virtual line breaks
	text = strings.Replace(text, `\n`, "\n", -1)

	// draw card rect
	fpdf.Rect(x, y, w, h, "D")

	// calc sizes
	cW, cH := w-(m*2), h-(m*2) // subtract two * margin
	lineHeight, lineCount := decreaseFont(fpdf, cW, cH, text)

	// loop and print lines with cell margin
	lX, lY := x+m, y+m+((cH-(lineHeight*lineCount))/2)
	fpdf.SetXY(lX, lY)
	for i, l := range cardLines(fpdf, w, text) {
		fpdf.CellFormat(cW, lineHeight, l, "", 0, "C", false, 0, "")
		fpdf.SetXY(lX, lY+float64(i+1)*lineHeight)
	}
}

// Print saves cards in a PDF file
func (p *PDF) Print(cards map[string]string) error {

	// create landscape PDF with unicode translation function
	fpdf := gofpdf.New("L", p.Card.Unit, p.Size, "")
	tr := fpdf.UnicodeTranslatorFromDescriptor("") // "" defaults to "cp1252"

	// calc sizes
	colWidth := float64(p.Card.Width)
	rowHeight := float64(p.Card.Height)
	textMargin := float64(p.Card.Margin)
	fontSize := float64(p.Font.Size)
	maxWidth, maxHeight := fpdf.GetPageSize()
	log.Printf("%s page %.2f x %.2f %s", p.Size, maxWidth, maxHeight, p.Card.Unit)

	// setup PDF
	//fpdf.AddFont(p.FontName, "", fmt.Sprintf("%s.json", p.FontName))
	fpdf.SetFont(p.Font.Name, "", fontSize)
	fpdf.SetAutoPageBreak(false, 0.0)
	fpdf.SetCellMargin(4.0)

	// calc counts
	cols := int(maxWidth / colWidth)
	rows := int(maxHeight / rowHeight)
	count := len(cards)
	pages := count / (cols * rows)
	if count%(cols*rows) != 0 {
		pages++ // add the rest
	}
	log.Printf("%d flashcard %s pages %d rows x %d cols", pages, p.Size, rows, cols)

	// calc margins
	colHeight := float64(rows * p.Card.Height)
	topMargin := (maxHeight - colHeight) / 2
	rowWidth := float64(cols * p.Card.Width)
	leftMargin := (maxWidth - rowWidth) / 2
	log.Printf("%.2f %s top and %.2f %s left margin", topMargin, p.Card.Unit, leftMargin, p.Card.Unit)

	printBack := func(back []string) {
		fpdf.AddPage()

		// reverse start position for the backside
		x, y := leftMargin, maxHeight-topMargin-rowHeight

		// iterate back
		for _, b := range back {

			// print back
			fpdf.SetFontSize(fontSize)
			printCard(fpdf, x, y, colWidth, rowHeight, textMargin, tr(b))
			x += colWidth // next col

			if x+colWidth+leftMargin > maxWidth { // next row
				x = leftMargin // reset
				y -= rowHeight // decrease
			}
		}
	}

	// add first page
	x, y := leftMargin, topMargin
	fpdf.AddPage()

	// iterate cards and remember backsides
	back := []string{}
	for front := range cards {
		back = append(back, cards[front])

		// print front
		fpdf.SetFontSize(fontSize)
		printCard(fpdf, x, y, colWidth, rowHeight, textMargin, tr(front))
		x += colWidth // next col

		// check next col width
		if x+colWidth+leftMargin > maxWidth { // next row
			x = leftMargin // reset
			y += rowHeight // increase

			if y+rowHeight+topMargin > maxHeight { // next page
				printBack(back)

				// reinit term loop
				back = []string{}
				fpdf.AddPage()
				x, y = leftMargin, topMargin
			}
		}
	}

	// iterate last back translation
	if count%(cols*rows) != 0 {
		printBack(back)
	}

	// write and close file
	err := fpdf.OutputFileAndClose(p.File)
	if err != nil {
		log.Println("Output FAILURE")
		return err
	}
	return nil
}
