package position

type Position struct {
	Index    int
	Line     int
	Column   int
	FileName string
	FileText string
}

func NewPosition(index int, line int, column int, fileName string, fileText string) Position {
	return Position{Index: index, Line: line, Column: column, FileName: fileName, FileText: fileText}
}

func (p *Position) Advance(currentChar *rune) {
	p.Index += 1
	p.Column += 1

	if currentChar != nil && *currentChar == '\n' {
		p.Line += 1
		p.Column = 0
	}
}

func (p Position) Copy() Position {
	return NewPosition(p.Index, p.Line, p.Column, p.FileName, p.FileText)
}
