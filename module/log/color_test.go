package glog

import (
	"fmt"
	"testing"
)

func TestColor(t *testing.T) {
	col := Color{1}
	if col.Value() != 1 {
		t.Errorf("col.Value() != 1, was %d", col.Value())
	}

	expected := "\u001b[31m"
	if col.String() != expected {
		t.Errorf("expected col.String() to be %q, was %q", expected, col.String())
	}

	expected = fmt.Sprintf("%shello%s", expected, ResetColor)
	if col.Color("hello") != expected {
		t.Errorf("expected col.Color() to be %q, was %q", expected, col.Color("hello"))
	}
}

func TestTextStyle(t *testing.T) {
	text := TextStyle{1, 22}

	expected := "\u001b[1m\u001b[22m"
	if text.String() != expected {
		t.Errorf("expected text.String() to be %q, was %q", expected, text.String())
	}

	expected = "\u001b[1mhello\u001b[22m"
	if text.TextStyle("hello") != expected {
		t.Errorf("expected text.TextStyle() to be %q, was %q", expected, text.String())
	}

	empty := TextStyle{}
	if empty.TextStyle("hello") != "hello" {
		t.Errorf("expected empty.TextStyle() to be %q, was %q",
			expected,
			empty.String())
	}
}

func TestStyle(t *testing.T) {
	colStyle := Red.NewStyle()
	expected := "\u001b[40m\u001b[31m"
	if colStyle.String() != expected {
		t.Errorf("expected colStyle.String() to be %q, was %q",
			expected,
			colStyle.String())
	}

	textStyle := Bold.NewStyle()
	expected = "\u001b[40m\u001b[30m\u001b[1mhello\u001b[22m\u001b[49m\u001b[39m"
	if textStyle.Style("hello") != expected {
		t.Errorf("expected textStyle.Style(\"hello\") to be %q, was %q",
			expected,
			textStyle.Style("hello"))
	}

	// reset it
	colStyle.Foreground(ResetColor)
	colStyle.Background(ResetColor)
	expected = "\u001b[49m\u001b[39m"
	if colStyle.String() != expected {
		t.Errorf("expected colStyle.String() to be %q, was %q",
			expected,
			colStyle.String())
	}

	builtStyle := colStyle.
		WithForeground(Red).
		WithBackground(Blue).
		WithTextStyle(Underline)
	expected = "\u001b[44m\u001b[31m\u001b[4mhello\u001b[24m\u001b[49m\u001b[39m"
	if builtStyle.Style("hello") != expected {
		t.Errorf("expected builtStyle.Style() to be %q, was %q",
			expected,
			builtStyle.Style("hello"))
	}
}

func TestExample(t *testing.T) {
	fmt.Println(Gray.Color("Writing in colors"))
	// You can just use colors
	fmt.Println(Red, "Writing in colors", Cyan, "is so much fun", Reset)
	fmt.Println(Magenta.Color("You can use colors to color specific phrases"))

	// You can just use text styles
	fmt.Println(Bold.TextStyle("We can have bold text"))
	fmt.Println(Underline.TextStyle("We can have underlined text"))
	fmt.Println(Bold, "But text styles don't work quite like colors :(")

	// Or you can use styles
	blueOnWhite := Blue.NewStyle().WithBackground(White)
	fmt.Printf("%s%s%s\n", blueOnWhite, "And they also have backgrounds!", Reset)
	fmt.Println(
		blueOnWhite.Style("You can style strings the same way you can color them!"))
	fmt.Println(
		blueOnWhite.WithTextStyle(Bold).
			Style("You can mix text styles with colors, too!"))

	// You can also easily make styling functions thanks to go's functional side
	lime := Green.NewStyle().
		WithBackground(Black).
		WithTextStyle(Bold).
		Style
	fmt.Println(lime("look at this cool lime text!"))
}
