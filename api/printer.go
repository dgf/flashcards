package api

// Printer that prints flash cards.
type Printer interface {

	// Print cards
	Print(cards map[string]string) error
}
