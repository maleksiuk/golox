package srccode

type sourceLocation struct {
	Start   int
	Current int
	Line    int
}

type Source struct {
	runes    []rune
	location sourceLocation
}

// NewSource creates a new Source based on the provided source code.
func NewSource(src string) Source {
	runes := []rune(src)
	location := sourceLocation{Line: 1}
	return Source{location: location, runes: runes}
}

// Len returns the length of the source code.
func (source *Source) Len() int {
	return len(source.runes)
}

// CurrentLine returns the line number that we are on.
func (source *Source) CurrentLine() int {
	return source.location.Line
}

// IncrementLine adds one to the line number.
func (source *Source) IncrementLine() {
	source.location.Line++
}

func (location *sourceLocation) atEnd(runes []rune) bool {
	return location.Current >= len(runes)
}

func (location *sourceLocation) beginNewLexeme() {
	location.Start = location.Current
}

// Substring returns a portion of the source based on the start and current positions.
func (source *Source) Substring(startOffset int, endOffset int) string {
	return string(source.runes[source.location.Start+startOffset : source.location.Current+endOffset])
}

// BeginNewLexeme adjusts our start location to the current location.
func (source *Source) BeginNewLexeme() {
	source.location.beginNewLexeme()
}

// AtEnd returns true if we have reached the end of the source code.
func (source *Source) AtEnd() bool {
	return source.location.atEnd(source.runes)
}

// Advance returns the current rune and then moves us on to the next rune.
func (source *Source) Advance() rune {
	r := source.currentRune()
	source.location.Current++
	return r
}

func (source *Source) currentRune() rune {
	return source.runes[source.location.Current]
}

func (source *Source) nextRune() rune {
	return source.runes[source.location.Current+1]
}

// Match advances us by one rune if our current rune matches the one passed in.
func (source *Source) Match(expected rune) bool {
	if source.AtEnd() {
		return false
	}

	if source.currentRune() != expected {
		return false
	}

	source.Advance()
	return true
}

// Peek returns the current rune without advancing.
func (source *Source) Peek() rune {
	if source.AtEnd() {
		return 0
	}

	return source.currentRune()
}

// PeekNext returns the next rune without advancing.
func (source *Source) PeekNext() rune {
	if source.location.Current+1 >= len(source.runes) {
		return 0
	}

	return source.nextRune()
}
