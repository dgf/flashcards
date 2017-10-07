package pdf

import "github.com/dgf/flashcards/config"

// Options of a flashcard PDF
type Options struct {
	Card config.Card
	Font config.Font
	Size string // A[0-6]
}
