package view_test

import (
	"testing"

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
