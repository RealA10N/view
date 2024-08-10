package view_test

import (
	"testing"
	"unicode"

	"alon.kr/x/view"
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
	v := view.NewView[int, uint]([]int{1, 2, 3, 2, 5})
	assert.EqualValues(t, 1, v.Index(2))
}

func TestIndexNotFound(t *testing.T) {
	v := view.NewView[int, uint]([]int{1, 2, 3, 4, 5})
	assert.EqualValues(t, v.Len(), v.Index(6))
}

func TestIndexFuncSimpleCase(t *testing.T) {
	v := view.NewView[int, uint]([]int{1337, -710, 2902}).Subview(1, 3)
	called := []int{}
	idx := v.IndexFunc(func(n int) bool { called = append(called, n); return n > 0 })
	assert.EqualValues(t, 1, idx)
	assert.Equal(t, []int{-710, 2902}, called)
}

func TestIndexFuncNotFound(t *testing.T) {
	data := []int{-1, 2, 3, -4}
	v := view.NewView[int, uint](data)
	called := []int{}
	idx := v.IndexFunc(func(n int) bool { called = append(called, n); return n > 10 })
	assert.EqualValues(t, len(data), idx)
	assert.Equal(t, data, called)
}

func TestContainsSimpleCase(t *testing.T) {
	data := []int{1, 2, 3}
	v := view.NewView[int, uint](data)
	assert.True(t, v.Contains(1))
	assert.False(t, v.Subview(1, 3).Contains(1))
}

func TestHasPrefixSimpleCase(t *testing.T) {
	v := view.NewView[int, uint]([]int{0, 1, 2, 3, 4, 5}).Subview(2, 5)
	prefix := view.NewView[int, uint]([]int{-1, 2, 3, 0}).Subview(1, 3)
	assert.True(t, v.HasPrefix(prefix))
}

func TestHasSuffixSimpleCase(t *testing.T) {
	v := view.NewView[int, uint]([]int{0, 1, 2, 3, 4, 5}).Subview(2, 5)
	suffix := view.NewView[int, uint]([]int{3, 4})
	assert.True(t, v.HasSuffix(suffix))
}

func TestLongestCommonPrefixSimpleCase(t *testing.T) {
	v := view.NewView[int, uint]([]int{0, 1, 2, 3, 4, 5}).Subview(1, 6)
	u := view.NewView[int, uint]([]int{1, 2, 4, 5})
	prefix := v.LongestCommonPrefix(u)
	assert.Equal(t, []int{1, 2}, prefix.Raw())
}

func TestLongestCommonPrefixMaxLength(t *testing.T) {
	v := view.NewView[int, uint]([]int{0, 1, 2, 3, 4, 5}).Subview(1, 6)
	u := view.NewView[int, uint]([]int{1, 2, 3, 4})
	prefix := v.LongestCommonPrefix(u)
	assert.EqualValues(t, []int{1, 2, 3, 4}, prefix.Raw())
}

func TestLongestCommonPrefixNoCommon(t *testing.T) {
	v := view.NewView[int, uint]([]int{1, 2, 3, 4, 5})
	u := view.NewView[int, uint]([]int{6, 7, 8, 9, 10})
	prefix := v.LongestCommonPrefix(u)
	assert.Equal(t, []int{}, prefix.Raw())
}

func TestLongestCommonSuffixSimpleCase(t *testing.T) {
	v := view.NewView[int, uint]([]int{1, 2, 3, 4, 5, 6}).Subview(0, 5)
	u := view.NewView[int, uint]([]int{2, 3, 4, 5})
	suffix := v.LongestCommonSuffix(u)
	assert.Equal(t, []int{2, 3, 4, 5}, suffix.Raw())
}

func TestLongestCommonSuffixMaxLength(t *testing.T) {
	v := view.NewView[int, uint]([]int{0, 1, 2, 3, 4, 5}).Subview(0, 5)
	u := view.NewView[int, uint]([]int{2, 3, 4})
	suffix := v.LongestCommonSuffix(u)
	assert.EqualValues(t, []int{2, 3, 4}, suffix.Raw())
}

func TestLongestCommonSuffixNoCommon(t *testing.T) {
	v := view.NewView[int, uint]([]int{1, 2, 3, 4, 5})
	u := view.NewView[int, uint]([]int{6, 7, 8, 9, 10})
	suffix := v.LongestCommonSuffix(u)
	assert.Equal(t, []int{}, suffix.Raw())
}

func TestMergeSimpleCase(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	v := view.NewView[int, uint](data)
	a := v.Subview(1, 3)
	b := v.Subview(5, 6)
	m := a.Merge(b)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, m.Raw())
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
