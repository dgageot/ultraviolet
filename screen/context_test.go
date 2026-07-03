package screen

import (
	"testing"
	"unicode/utf8"
)

func TestDrawStringInvalidUTF8(t *testing.T) {
	scr := newMockScreen(20, 1)
	ctx := NewContext(scr)

	// Text cropped mid-emoji, leaving a trailing invalid UTF-8 sequence.
	ctx.DrawString("hi \xf0\x9f\xab", 0, 0)
	ctx.DrawString("│", 6, 0)

	for x := range 20 {
		c := scr.CellAt(x, 0)
		if !utf8.ValidString(c.Content) {
			t.Errorf("cell at %d contains invalid UTF-8: %q", x, c.Content)
		}
	}

	if got := scr.CellAt(3, 0).Content; got != "\uFFFD" {
		t.Errorf("cell at 3 = %q, want %q", got, "\uFFFD")
	}
	if got := scr.CellAt(6, 0).Content; got != "│" {
		t.Errorf("cell at 6 = %q, want %q", got, "│")
	}
}

func TestWriteInvalidUTF8(t *testing.T) {
	scr := newMockScreen(20, 1)
	ctx := NewContext(scr)

	n, err := ctx.Write([]byte("a\xffb"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Errorf("Write returned %d, want 3", n)
	}

	want := []string{"a", "\uFFFD", "b"}
	for x, w := range want {
		if got := scr.CellAt(x, 0).Content; got != w {
			t.Errorf("cell at %d = %q, want %q", x, got, w)
		}
	}
}
