package view_test

import (
	"testing"
	"unicode"

	"github.com/RealA10N/view"
	"github.com/stretchr/testify/assert"
)

func TestSimpleSubview(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	v := view.NewView[int, uint](data).Subview(1, 4)
	assert.Equal(t, []int{1, 2, 3}, v.Raw())
}

func TestOutOfBoundsSubview(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	v := view.NewView[int, uint](data).Subview(10, 13)
	assert.Equal(t, []int{}, v.Raw())
}

func TestReversedBoundsSubview(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	v := view.NewView[int, uint](data).Subview(2, 1)
	assert.Equal(t, []int{}, v.Raw())
}

func TestDetach(t *testing.T) {
	v := view.NewView[int, uint]([]int{1, 2, 3})
	unmanaged, _ := v.Detach()
	assert.EqualValues(t, 0, unmanaged.Start)
	assert.EqualValues(t, 3, unmanaged.End)
}

func TestFields(t *testing.T) {
	input := []rune("  foo1;bar2,baz3...")
	v := view.NewView[rune, uint](input)
	fields := v.Fields(func(c rune) bool { return !unicode.IsLetter(c) && !unicode.IsNumber(c) })

	got := [][]rune{}
	for _, view := range fields {
		got = append(got, view.Raw())
	}

	expected := [][]rune{[]rune("foo1"), []rune("bar2"), []rune("baz3")}
	assert.Equal(t, expected, got)
}
