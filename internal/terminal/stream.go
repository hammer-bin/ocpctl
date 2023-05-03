package terminal

import "os"

const defaultColumns int = 78
const defaultIsTerminal bool = false

type OutputStream struct {
	File *os.File

	isTerminal func(*os.File) bool
	getColumns func(*os.File) int
}

func (s *OutputStream) Columns() int {
	if s.getColumns == nil {
		return defaultColumns
	}
	return s.getColumns(s.File)
}

func (s *OutputStream) IsTerminal() bool {
	if s.isTerminal == nil {
		return defaultIsTerminal
	}
	return s.isTerminal(s.File)
}

type InputStream struct {
	File *os.File

	isTerminal func(*os.File) bool
}

func (s *InputStream) IsTerminal() bool {
	if s.isTerminal == nil {
		return defaultIsTerminal
	}
	return s.isTerminal(s.File)
}