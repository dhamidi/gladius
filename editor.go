package gladius

import "strings"

// Editor is a buffer with an attached cursor that represents the
// current insert position in the buffer.
type Editor struct {
	cursor               int64
	buffer               *Buffer
	cachedBufferContents string
	cacheIsDirty         bool
}

// NewEditor creates an editor for a buffer with the provided contents.
func NewEditor(contents string) *Editor {
	return &Editor{
		buffer:               NewBufferString(contents),
		cacheIsDirty:         false,
		cachedBufferContents: contents,
	}
}

// Position returns the current cursor position of the editor.
func (ed *Editor) Position() int64 {
	return ed.cursor
}

// Backward moves the cursors backward by n charactrs.  If the cursors
// would go beyond the beginning of the buffer, the cursor position is
// instead set to the beginning of the buffer.
func (ed *Editor) Backward(n int64) *Editor {
	newPosition := ed.cursor - n
	if newPosition < 0 {
		ed.cursor = 0
		return ed
	}
	ed.cursor = newPosition
	return ed
}

// Forward moves the cursor forward by n characters.  If the cursor
// would reach the end of the buffer, it is positioned after the last
// character in the buffer.
func (ed *Editor) Forward(n int64) *Editor {
	if l := ed.buffer.Len(); ed.cursor+n >= l {
		ed.cursor = l
		return ed
	}
	ed.cursor += n
	return ed
}

// Insert inserts the given string at the current cursor position and
// advances the cursor by the length of the string.
func (ed *Editor) Insert(text string) *Editor {
	ed.buffer.Insert(ed.cursor, text)
	ed.cursor += int64(len(text))
	ed.cacheIsDirty = true
	return ed
}

// String returns the contents of the buffer managed by this editor.
func (ed *Editor) String() string {
	if ed.cacheIsDirty {
		ed.cachedBufferContents = ed.buffer.String()
		ed.cacheIsDirty = false
	}
	return ed.cachedBufferContents
}

// Before returns the contents of the buffer before the cursor.
func (ed *Editor) Before() string {
	return ed.String()[0:ed.cursor]
}

// After returns the contents of the buffer after the cursor.
func (ed *Editor) After() string {
	return ed.String()[ed.cursor:]
}

// Delete removes n characters after the cursor.
func (ed *Editor) Delete(n int64) *Editor {
	ed.buffer.Delete(ed.cursor, n)
	ed.cacheIsDirty = true
	return ed
}

// BeginningOfBuffer sets the cursor position to the beginning of the underlying buffer.
func (ed *Editor) BeginningOfBuffer() *Editor {
	ed.cursor = 0
	return ed
}

// EndOfBuffer sets the cursor position to the end of the underlying buffer.
func (ed *Editor) EndOfBuffer() *Editor {
	ed.cursor = ed.buffer.Len()
	return ed
}

// BeginningOfLine sets the cursor after the first newline character before the cursor.
func (ed *Editor) BeginningOfLine() *Editor {
	newPosition := strings.LastIndex(ed.Before(), "\n")
	if newPosition > -1 {
		ed.cursor = int64(newPosition) + int64(1)
	} else {
		ed.BeginningOfBuffer()
	}
	return ed
}

// EndOfLine sets the cursor to the first newline character after the cursor.
func (ed *Editor) EndOfLine() *Editor {
	offsetFromCursor := strings.Index(ed.After(), "\n")
	if offsetFromCursor > -1 {
		ed.cursor = ed.cursor + int64(offsetFromCursor)
		return ed
	}
	return ed.EndOfBuffer()
}

// ForwardLine moves the cursor forward by n lines.
func (ed *Editor) ForwardLine(n int) *Editor {
	for ; n > 0; n-- {
		ed.EndOfLine().Forward(1)
	}
	return ed
}

// BackwardLine moves the cursor backward by n lines.
func (ed *Editor) BackwardLine(n int) *Editor {
	for ; n > 0; n-- {
		ed.BeginningOfLine().Backward(1).BeginningOfLine()
	}
	return ed
}
