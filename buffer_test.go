package gladius

import "testing"

func TestBuffer_InsertTextAtEnd(t *testing.T) {
	buffer := NewBufferString("end")
	buffer.Insert(3, "hello")

	if exp, act := "endhello", buffer.String(); exp != act {
		t.Fatalf("Expected %q, got %q", exp, act)
	}
}

func TestBuffer_InsertTextAtBeginning(t *testing.T) {
	buffer := NewBufferString("llo")
	buffer.Insert(0, "he")

	if exp, act := "hello", buffer.String(); exp != act {
		t.Fatalf("Expected %q, got %q", exp, act)
	}
}

func TestBuffer_InsertTextInTheMiddle(t *testing.T) {
	buffer := NewBufferString("hello world")
	buffer.Insert(6, "cruel ")
	if exp, act := "hello cruel world", buffer.String(); exp != act {
		t.Fatalf("Expected %q, got %q", exp, act)
	}
}

func TestBuffer_MultipleInsertsInTheMiddle(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Insert(1, "1")
	buffer.Insert(3, "2")

	if exp, act := "a1b2c", buffer.String(); exp != act {
		t.Fatalf("Expected %q, got %q", exp, act)
	}
}

func TestBuffer_MultipleInsertsOfDifferentTypes(t *testing.T) {
	buffer := NewBufferString("abc")
	buffer.Insert(0, "1234")
	buffer.Insert(int64(len("abc")+len("1234")), "5678")
	buffer.Insert(6, "!")
	t.Logf("Buffer piece table:\n%s\n", buffer.Inspect())
	if exp, act := "1234ab!c5678", buffer.String(); exp != act {
		t.Fatalf("Expected %q, got %q", exp, act)
	}
}
