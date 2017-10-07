package config

// Font configuration
type Font struct {
	Name  string // TODO add dynamic font loader
	Size  int    // font size (points, not units)
	Style string // [BIU]
}
