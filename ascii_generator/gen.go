package ascii_generator

type Cell struct {
	Ch    rune
	Color string
}

type AsciiArt struct {
	Canvas [][]Cell
	Width  int
	Height int
}

type AsciiArtInterface interface {
	GenerateAsciiArt(int, int) AsciiArt
	StringArray() []string
}
