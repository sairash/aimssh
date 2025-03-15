package ascii_generator

type AsciiArt interface {
	Width() int
	Height() int
	Next(int) bool
	NextAndString(int) string
}
