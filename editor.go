package gladius

// Editor is a buffer with an attached cursor that represents the
// current insert position in the buffer.
type Editor struct {
	cursor int64
	buffer *Buffer
}

// NewEditor creates an editor for a buffer with the provided contents.
func NewEditor(contents string) *Editor {
	return &Editor{
		buffer: NewBufferString(contents),
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
	return ed
}

// String returns the contents of the buffer managed by this editor.
func (ed *Editor) String() string {
	return ed.buffer.String()
}

// Delete removes n characters after the cursor.
func (ed *Editor) Delete(n int64) *Editor {
	ed.buffer.Delete(ed.cursor, n)
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
	newPosition := ed.buffer.FindBackwards(ed.cursor, '\n')
	if newPosition > -1 {
		ed.cursor = newPosition + 1
	}
	return ed
}

// BeginningOfLine sets the cursor to the first newline character after the cursor.
func (ed *Editor) EndOfLine() *Editor {
	newPosition := ed.buffer.FindForwards(ed.cursor, '\n')
	if newPosition > -1 {
		ed.cursor = newPosition
		return ed
	}
	return ed.EndOfBuffer()
}
