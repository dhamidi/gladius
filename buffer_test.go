package gladius

import (
	"fmt"
	"math/rand"
	"testing"
)

func assertBuffer(t *testing.T, b *Buffer, expected string) {
	t.Helper()
	t.Logf("Buffer piece table:\n%s\n", b.Inspect())
	if contents := b.String(); contents != expected {
		t.Fatalf("Expected %q, got %q", expected, contents)
	}
}

func TestBuffer_InsertTextAtEnd(t *testing.T) {
	buffer := NewBufferString("end")
	buffer.Insert(3, "hello")

	assertBuffer(t, buffer, "endhello")
}

func TestBuffer_InsertTextAtBeginning(t *testing.T) {
	buffer := NewBufferString("llo")
	buffer.Insert(0, "he")

	assertBuffer(t, buffer, "hello")
}

func TestBuffer_InsertTextInTheMiddle(t *testing.T) {
	buffer := NewBufferString("hello world")
	buffer.Insert(6, "cruel ")
	assertBuffer(t, buffer, "hello cruel world")
}

func TestBuffer_MultipleInsertsInTheMiddle(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Insert(1, "1")
	buffer.Insert(3, "2")

	assertBuffer(t, buffer, "a1b2c")
}

func TestBuffer_MultipleInsertsOfDifferentTypes(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Insert(0, "1234")
	buffer.Insert(int64(len("abc")+len("1234")), "5678")
	buffer.Insert(6, "!")

	assertBuffer(t, buffer, "1234ab!c5678")
}

func TestBuffer_DeleteRemovesTextAtBeginning(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Delete(0, 1)

	assertBuffer(t, buffer, "bc")
}

func TestBuffer_DeleteRemovesTextInTheMiddle(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Delete(1, 1)

	assertBuffer(t, buffer, "ac")
}

func TestBuffer_DeleteRemovesTextAtEnd(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Delete(2, 1)

	assertBuffer(t, buffer, "ab")
}

func TestBuffer_MultipleMixedDeletes(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Delete(1, 1)
	buffer.Delete(0, 2)
	assertBuffer(t, buffer, "")
}

func TestBuffer_MultipleDeletesStartingInTheMiddle(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Delete(1, 1)
	buffer.Delete(1, 1)

	assertBuffer(t, buffer, "a")
}

func TestBuffer_MixedInsertsAndDeletes(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Delete(1, 1) // ac
	assertBuffer(t, buffer, "ac")
	buffer.Insert(1, "def") // adefc
	assertBuffer(t, buffer, "adefc")
	buffer.Delete(3, 2) // ade
	assertBuffer(t, buffer, "ade")
	buffer.Insert(0, "casc") // cascade

	assertBuffer(t, buffer, "cascade")
}

func TestBuffer_FindBackwards_returns_offset_of_character_before_given_position(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Insert(2, "def")
	locationOfB := buffer.FindBackwards(3, 'b')
	assertEqual(t, locationOfB, int64(1))
}

func TestBuffer_FindBackwards_returns_minus_one_if_character_can_not_be_found(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Insert(2, "def")
	locationOfX := buffer.FindBackwards(3, 'x')
	assertEqual(t, locationOfX, int64(-1))
}

func TestBuffer_FindForwards_returns_offset_of_character_after_given_position(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Insert(2, "def")
	locationOfF := buffer.FindForwards(3, 'f')
	assertEqual(t, locationOfF, int64(4))
}

func TestBuffer_FindForwards_returns_minus_one_if_character_can_not_be_found(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Insert(2, "def")
	locationOfX := buffer.FindForwards(3, 'x')
	assertEqual(t, locationOfX, int64(-1))
}

func BenchmarkSequentialInsertAtEnd(b *testing.B) {
	buffer := NewBufferString("")
	position := int64(0)
	for i := 0; i < b.N; i++ {
		text := fmt.Sprintf("%d", b.N)
		buffer.Insert(position, text)
		position = position + int64(len(text))
	}
}

func BenchmarkSequentialInsertAtBeginning(b *testing.B) {
	buffer := NewBufferString("")
	for i := 0; i < b.N; i++ {
		text := fmt.Sprintf("%d", b.N)
		buffer.Insert(0, text)
	}
}

func BenchmarkRandomInsert(b *testing.B) {
	buffer := NewBufferString("")
	for i := 0; i < b.N; i++ {
		position := int64(0)
		if len := buffer.Len(); len > 0 {
			position = rand.Int63n(len)
		}
		buffer.Insert(position, fmt.Sprintf("%d", position))
	}
}
