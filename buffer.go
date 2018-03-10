package gladius

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Buffer is a text buffer supporting simple edit operations at random
// positions in the buffer's text.
//
// This text buffer is implemented using a piece table.
type Buffer struct {
	base   string
	add    *bytes.Buffer
	pieces []*piece
}

// pieceBufferType is an enum describing to which buffer a piece in the piece table pertains.
type pieceBufferType int

const (
	pieceBufferBase = pieceBufferType(1 << iota)
	pieceBufferAdd  = pieceBufferType(1 << iota)
)

func (t pieceBufferType) String() string {
	switch t {
	case pieceBufferBase:
		return "base"
	case pieceBufferAdd:
		return "add"
	default:
		panic(fmt.Sprintf("unknown piece buffer type: %d", int(t)))
	}
}

// piece represents a span of text, either from the base or add buffer.
type piece struct {
	buffer pieceBufferType //
	offset int64           // the start of the span in the given buffer
	length int64           // the length of the span in the given buffer
}

// split splits a piece into two separate pieces at position i.
func (p *piece) split(i int64) []*piece {
	before := &piece{
		buffer: p.buffer,
		offset: p.offset,
		length: i,
	}
	after := &piece{
		buffer: p.buffer,
		offset: before.offset + before.length,
		length: p.length - before.length,
	}
	result := []*piece{before}
	if after.length > 0 {
		result = append(result, after)
	}
	return result
}

// text returns the text pointed at by this piece for the given buffer
func (p *piece) text(b *Buffer) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "DEBUG: %#v\n", err)
		}
	}()
	add := b.add.String()
	text := ""
	switch p.buffer {
	case pieceBufferAdd:
		text = add[p.offset : p.offset+p.length]
	case pieceBufferBase:
		text = b.base[p.offset : p.offset+p.length]
	}
	return text
}

// NewBufferString creates a new buffer using str as the initial text.
func NewBufferString(str string) *Buffer {
	return &Buffer{
		base: str,
		add:  bytes.NewBufferString(""),
		pieces: []*piece{
			{
				buffer: pieceBufferBase,
				offset: 0,
				length: int64(len(str)),
			},
		},
	}
}

// pieceAt returns the piece that contains location loc and the index
// of that piece in the piece table.  If no such piece is found, nil
// and -1 is returned.
func (b *Buffer) pieceAt(loc int64) (*piece, int64, int) {
	length := int64(0)
	for i, piece := range b.pieces {
		if loc >= length && loc <= length+piece.length {
			return piece, length, i
		}
		length = length + piece.length
	}

	return nil, loc, -1
}

// Inspect returns the internal piece table of the buffer as a string.
func (b *Buffer) Inspect() string {
	out := bytes.NewBufferString("")
	dumpPieces(out, b.pieces, b)
	return out.String()
}

// Insert inserts text at loc in the buffer.
func (b *Buffer) Insert(loc int64, text string) *Buffer {
	currentPiece, offset, listIndex := b.pieceAt(loc)
	newPiece := &piece{
		buffer: pieceBufferAdd,
		offset: int64(b.add.Len()),
		length: int64(len(text)),
	}
	fmt.Fprintf(b.add, "%s", text)
	if currentPiece != nil {
		split := currentPiece.split(loc - offset)
		before, after := split[0], split[1:]
		rest := make([]*piece, len(b.pieces[listIndex+1:]))
		copy(rest, b.pieces[listIndex+1:])
		if listIndex == 0 {
			b.pieces = append(
				[]*piece{
					before,
					newPiece,
				},
				after...,
			)
			b.pieces = append(b.pieces, rest...)
		} else {

			b.pieces = append(b.pieces[0:listIndex], before, newPiece)
			b.pieces = append(b.pieces, after...)
			b.pieces = append(b.pieces, rest...)
		}
		return b
	}

	b.pieces = append(b.pieces, newPiece)

	return b
}

func dumpPieces(out io.Writer, pieces []*piece, b *Buffer) {
	for i, piece := range pieces {
		fmt.Fprintf(out, "%3d %5s %3d %3d %q\n", i, piece.buffer, piece.offset, piece.length, piece.text(b))
	}
}

// Delete removes n characters at loc in the buffer.
func (b *Buffer) Delete(loc int64, n int64) *Buffer {
	endLoc := loc + n
	beginPiece, beginOffset, beginIndex := b.pieceAt(loc)
	endPiece, endOffset, endIndex := b.pieceAt(endLoc)
	beginSplit := beginPiece.split(loc - beginOffset)
	endSplit := endPiece.split(endLoc - endOffset)
	result := []*piece{beginSplit[0]}
	if len(endSplit) > 1 {
		result = append(result, endSplit[1])
	}
	if endIndex < len(b.pieces) {
		result = append(result, b.pieces[endIndex+1:]...)
	}
	if beginIndex == 0 {
		b.pieces = result
	} else {
		b.pieces = append(
			b.pieces[0:beginIndex],
			result...,
		)
	}

	return b
}

// String returns the current contents of the buffer as a single string.
func (b *Buffer) String() string {
	out := bytes.NewBufferString("")
	add := b.add.String()
	for _, piece := range b.pieces {
		text := ""
		switch piece.buffer {
		case pieceBufferAdd:
			text = add[piece.offset : piece.offset+piece.length]
		case pieceBufferBase:
			text = b.base[piece.offset : piece.offset+piece.length]
		}
		fmt.Fprintf(out, "%s", text)
	}
	return out.String()
}
