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

func TestEqualSimpleCase(t *testing.T) {
	a := view.NewView[rune, uint]([]rune("gila bisa tomer natanel")).Subview(5, 15)
	b := view.NewView[rune, uint]([]rune("bisa tomer lorem ipsum")).Subview(0, 10)
	assert.True(t, a.Equal(b))
	assert.True(t, b.Equal(a))
}

func TestEqualSubstrings(t *testing.T) {
	a := view.NewView[rune, uint]([]rune("gila bisa"))
	b := a.Subview(0, 8)
	assert.False(t, a.Equal(b))
	assert.False(t, b.Equal(a))
}

func TestIndexSimpleCase(t *testing.T) {
	v := view.NewView[int, uint]([]int{1337, -710, 2902}).Subview(1, 3)
	called := []int{}
	idx, err := v.Index(func(n int) bool { called = append(called, n); return n > 0 })
	assert.NoError(t, err)
	assert.EqualValues(t, 1, idx)
	assert.Equal(t, []int{-710, 2902}, called)
}

func TestIndexNotFound(t *testing.T) {
	data := []int{-1, 2, 3, -4}
	v := view.NewView[int, uint](data)
	called := []int{}
	_, err := v.Index(func(n int) bool { called = append(called, n); return n > 10 })
	assert.Error(t, err)
	assert.Equal(t, data, called)
}

func TestPartitionSimpleCase(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	v := view.NewView[int, uint](data)
	a, b := v.Partition(2)
	assert.Equal(t, []int{1, 2}, a.Raw())
	assert.Equal(t, []int{3, 4, 5}, b.Raw())
}

func TestPartitionZeroIndex(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	v := view.NewView[int, uint](data)
	a, b := v.Partition(0)
	assert.Equal(t, []int{}, a.Raw())
	assert.Equal(t, []int{1, 2, 3, 4, 5}, b.Raw())
}

func TestPartitionOutOfBoundsIndex(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	v := view.NewView[int, uint](data)
	a, b := v.Partition(10)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a.Raw())
	assert.Equal(t, []int{}, b.Raw())
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
