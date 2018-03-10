package gladius

import "testing"

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
