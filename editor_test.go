package gladius

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		return
	}

	t.Fatalf("Expected %v to equal %v", a, b)
}

func TestEditor_Position_in_an_empty_buffer_returns_0_initially(t *testing.T) {
	editor := NewEditor("")
	assertEqual(t, editor.Position(), int64(0))
}

func TestEditor_Forward_increases_position(t *testing.T) {
	editor := NewEditor("123")
	editor.Forward(1).Forward(1)
	assertEqual(t, editor.Position(), int64(2))
}

func TestEditor_Forward_in_an_empty_editor_does_not_alter_position(t *testing.T) {
	editor := NewEditor("")
	editor.Forward(10)
	assertEqual(t, editor.Position(), int64(0))
}

func TestEditor_Forward_stops_at_the_end_of_the_buffer(t *testing.T) {
	editor := NewEditor("asdf")
	editor.Insert("ghjkl")
	editor.Forward(20)
	assertEqual(t, editor.Position(), int64(len("asdfghjkl")))
}

func TestEditor_Backward_decreases_position(t *testing.T) {
	editor := NewEditor("123")
	editor.Forward(3).Backward(1)
	assertEqual(t, editor.Position(), int64(2))
}

func TestEditor_Backward_in_an_empty_editor_does_not_alter_position(t *testing.T) {
	editor := NewEditor("")
	editor.Backward(10)
	assertEqual(t, editor.Position(), int64(0))
}

func TestEditor_Backward_stops_at_the_beginning_of_the_buffer(t *testing.T) {
	editor := NewEditor("asdf")
	editor.Insert("ghjkl")
	editor.Backward(20)
	assertEqual(t, editor.Position(), int64(0))
}

func TestEditor_Insert_adds_text_at_the_current_position(t *testing.T) {
	editor := NewEditor("123")
	editor.Forward(1)
	editor.Insert("asdf")
	assertEqual(t, editor.String(), "1asdf23")
}

func TestEditor_Insert_advances_the_current_cursor_position(t *testing.T) {
	editor := NewEditor("123")
	editor.Forward(1)
	editor.Insert("asdf")
	assertEqual(t, editor.Position(), int64(1+len("asdf")))
}

func TestEditor_Delete_removes_text_after_the_cursor(t *testing.T) {
	editor := NewEditor("123")
	editor.Delete(2)
	assertEqual(t, editor.String(), "3")
}

func TestEditor_BeginningOfBuffer_sets_position_to_0(t *testing.T) {
	editor := NewEditor("123")
	editor.Forward(2).BeginningOfBuffer()
	assertEqual(t, editor.Position(), int64(0))
}

func TestEditor_EndOfBuffer_sets_position_to_length_of_inserted_text(t *testing.T) {
	editor := NewEditor("123")
	editor.Insert("asdf").EndOfBuffer()
	assertEqual(t, editor.Position(), int64(len("123")+len("asdf")))
}

func TestEditor_BeginningOfLine_sets_position_to_first_newline_before_cursor(t *testing.T) {
	editor := NewEditor("hello\n").EndOfBuffer()
	editor.Insert("orld").BeginningOfLine().Insert("w")
	assertEqual(t, editor.String(), "hello\nworld")
	assertEqual(t, editor.Position(), int64(len("hello\n")+1))
}

func TestEditor_BeginningOfLine_moves_to_beginning_of_buffer_if_no_newline_is_before_the_cursor(t *testing.T) {
	editor := NewEditor("hello").EndOfBuffer()
	editor.Insert("orld").BeginningOfLine().Insert("w")
	assertEqual(t, editor.String(), "whelloorld")
	assertEqual(t, editor.Position(), int64(1))
}

func TestEditor_EndOfLine_sets_position_to_first_newline_after_cursor(t *testing.T) {
	editor := NewEditor("hello\n").BeginningOfBuffer()
	editor.EndOfLine().Insert(" world")
	assertEqual(t, editor.String(), "hello world\n")
	assertEqual(t, editor.Position(), int64(len("hello world")))
}

func TestEditor_EndOfLine_moves_to_end_of_buffer_if_no_newline_occurs_after_the_cursor(t *testing.T) {
	editor := NewEditor("hello").EndOfLine()
	assertEqual(t, editor.Position(), int64(len("hello")))
}

func TestEditor_ForwardLine_moves_the_cursor_to_the_beginning_of_the_next_line(t *testing.T) {
	editor := NewEditor("").
		Insert("hello\n").
		Insert("world\n").
		BeginningOfBuffer().
		ForwardLine(1).
		Insert(", ")
	assertEqual(t, editor.String(), "hello\n, world\n")
	assertEqual(t, editor.Position(), int64(len("hello\n, ")))
}

func TestEditor_ForwardLine_moves_the_cursor_through_multiple_lines(t *testing.T) {
	editor := NewEditor("").
		Insert("hello\n").
		Insert("world\n").
		BeginningOfBuffer().
		ForwardLine(2).
		Insert("!")
	assertEqual(t, editor.String(), "hello\nworld\n!")
	assertEqual(t, editor.Position(), int64(len("hello\nworld\n!")))
}

func TestEditor_BackwardLine_moves_the_cursor_to_the_beginning_of_the_previous(t *testing.T) {
	editor := NewEditor("").
		Insert("hello\n").
		Insert("world\n").
		EndOfBuffer().
		BackwardLine(1).
		Insert(", ")
	assertEqual(t, editor.String(), "hello\n, world\n")
	assertEqual(t, editor.Position(), int64(len("hello\n, ")))
}

func TestEditor_BackwardLine_moves_the_cursor_through_multiple_lines(t *testing.T) {
	editor := NewEditor("").
		Insert("hello\n").
		Insert("world\n").
		EndOfBuffer().
		BackwardLine(2).
		Insert("!")
	assertEqual(t, editor.String(), "!hello\nworld\n")
	assertEqual(t, editor.Position(), int64(len("!")))
}
