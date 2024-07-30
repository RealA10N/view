package view

import (
	"golang.org/x/exp/constraints"
)

// The most basic slice view type.
// Internally, it contains a pointer to a heap allocated slice of type []T,
// and Start & End indecies of type Offset.
type BasicView[T comparable, Offset constraints.Unsigned] struct {
	Start, End Offset
	data       *[]T
}

// Create a new basic view from an already existing slice.
// The view initially spans over the whole slice.
func NewBasicView[T comparable, Offset constraints.Unsigned](data []T) BasicView[T, Offset] {
	return BasicView[T, Offset]{
		Start: 0,
		End:   Offset(len(data)),
		data:  &data,
	}
}

// Create a new view from an already existing slice.
// The view initially spans over the whole slice.
func NewView[T comparable](data []T) BasicView[T, uint] {
	return NewBasicView[T, uint](data)
}

// Convert the view into a slice and return a copy of the viewed slice only.
func (v BasicView[T, Offset]) Raw() []T {
	return (*v.data)[v.Start:v.End]
}

// Returns the size of the view slice.
func (v BasicView[T, Offset]) Len() uint {
	return uint(v.End - v.Start)
}

// Returns the item at the provided index, relative to the view bounds.
func (v BasicView[T, Offset]) At(index Offset) *T {
	index += v.Start
	if index >= v.End {
		return nil
	}
	return &(*v.data)[index]
}

// Iterate over all values in the view (rangefunc).
func (v BasicView[T, Offset]) Range(yield func(T) bool) {
	for i := v.Start; i < v.End; i++ {
		if !yield((*v.data)[i]) {
			return
		}
	}
}

// Iterate over all values in the view (rangefunc).
// Additionally, provides the iteration index as the first yield argument,
// where the index is relative to the view start.
func (v BasicView[T, Offset]) Range2(yield func(Offset, T) bool) {
	var sliceIndex Offset = v.Start
	var viewIndex Offset = 0

	for sliceIndex < v.End {
		if !yield(viewIndex, (*v.data)[sliceIndex]) {
			return
		}

		sliceIndex++
		viewIndex++
	}
}

// Compares two views and returns true if and only their lengths are equal
// and the corresponding items in both views are equal.
func (v BasicView[T, Offset]) Equals(o BasicView[T, Offset]) bool {
	if v.Len() != o.Len() {
		return false
	}

	ret := true
	yield := func(index Offset, item T) bool {
		if *o.At(index) != item {
			ret = false
			return false
		}
		return true
	}

	v.Range2(yield)
	return ret
}
