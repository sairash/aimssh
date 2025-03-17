package ascii_generator

type AsciiArt interface {
	Width() int
	Height() int
	NextAndString(int) string
}
