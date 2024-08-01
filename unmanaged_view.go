package view

import (
	"errors"

	"golang.org/x/exp/constraints"
)

// The most basic slice view type.
// Internally, it contains a pointer to a heap allocated slice of type []T,
// and Start & End indecies of type Offset.
type UnmanagedView[T comparable, Offset constraints.Unsigned] struct {
	Start, End Offset
}

type ViewContext[T any] []T

// Create a new basic view from an already existing slice.
// The view initially spans over the whole slice.
func NewUnmanagedView[T comparable, Offset constraints.Unsigned](data []T) (
	UnmanagedView[T, Offset], ViewContext[T],
) {
	view := UnmanagedView[T, Offset]{
		Start: 0,
		End:   Offset(len(data)),
	}

	ctx := make(ViewContext[T], len(data))
	copy(ctx, data)

	return view, ctx
}

func (v UnmanagedView[T, Offset]) Raw(ctx ViewContext[T]) []T {
	return ctx[v.Start:v.End]
}

// Returns the size of the view slice.
func (v UnmanagedView[T, Offset]) Len() Offset {
	return v.End - v.Start
}

// Returns the item at the provided index, relative to the view bounds.
func (v UnmanagedView[T, Offset]) At(ctx ViewContext[T], index Offset) (T, error) {
	index += v.Start
	if index >= v.End {
		var t T
		return t, errors.New("index out of view bounds")
	}
	return ctx[index], nil
}

func (v UnmanagedView[T, Offset]) Subview(start, end Offset) (subview UnmanagedView[T, Offset]) {
	if start > end {
		// provided start index greater than end index
		return
	}

	if v.Start+end > v.End {
		// provided end index is out of current view bound
		return
	}

	subview = UnmanagedView[T, Offset]{
		Start: v.Start + start,
		End:   v.Start + end,
	}

	return
}

// Iterate over all values in the view (rangefunc).
func (v UnmanagedView[T, Offset]) Range(ctx ViewContext[T]) func(func(T) bool) {
	return func(yield func(T) bool) {
		for i := v.Start; i < v.End; i++ {
			if !yield(ctx[i]) {
				return
			}
		}
	}
}

// Iterate over all values in the view (rangefunc).
// Additionally, provides the iteration index as the first yield argument,
// where the index is relative to the view start.
func (v UnmanagedView[T, Offset]) Range2(ctx ViewContext[T]) func(func(Offset, T) bool) {
	return func(yield func(Offset, T) bool) {
		var sliceIndex Offset = v.Start
		var viewIndex Offset = 0

		for sliceIndex < v.End {
			if !yield(viewIndex, ctx[sliceIndex]) {
				return
			}

			sliceIndex++
			viewIndex++
		}
	}
}
