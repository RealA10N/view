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

// Create a new (unmanaged) view from an already existing slice.
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

func (v UnmanagedView[T, Offset]) Attach(ctx ViewContext[T]) View[T, Offset] {
	return View[T, Offset]{
		unmanaged: v,
		ctx:       ctx,
	}
}

// Returns the raw underlying slice that the view is bound to.
func (v UnmanagedView[T, Offset]) Raw(ctx ViewContext[T]) []T {
	return ctx[v.Start:v.End]
}

// Returns the size of the view slice.
func (v UnmanagedView[T, Offset]) Len() Offset {
	return v.End - v.Start
}

// Returns the item at the provided index, relative to the view bounds.
// If the provided index goes out of the view bounds, an error is returned,
// with an undefined value.
func (v UnmanagedView[T, Offset]) At(ctx ViewContext[T], index Offset) (T, error) {
	index += v.Start
	if index >= v.End {
		var t T
		return t, errors.New("index out of view bounds")
	}
	return ctx[index], nil
}

// Returns the item at the provided index, relative to the view bounds.
// Does not check view bounds, and if the provided index is greater than the
// view length, the function might panic or return an undefined value.
func (v UnmanagedView[T, Offset]) AtUnsafe(ctx ViewContext[T], index Offset) T {
	return ctx[v.Start+index]
}

// Return a subview of the current view.
// Start and End indecies are relative to the current view bounds,
// i.e. v.Subview(0, v.Len()) will return a subview that equals to the current one.
func (v UnmanagedView[T, Offset]) Subview(start, end Offset) UnmanagedView[T, Offset] {
	len := v.Len()
	if end > len {
		end = len
	}

	if start > end {
		start = end
	}

	return UnmanagedView[T, Offset]{
		Start: v.Start + start,
		End:   v.Start + end,
	}
}

// returns true if the underlying views are identical in their content.
//
// The first argument is the context of the this view.
// Then comes the other unmanaged view, and the third argument is the context
// of the other unmanaged view which was provided in the second argument.
func (v UnmanagedView[T, Offset]) Equal(
	vctx ViewContext[T], o UnmanagedView[T, Offset], octx ViewContext[T],
) bool {
	if v.Len() != o.Len() {
		return false
	}

	len := v.Len()
	var i Offset = 0
	for ; i < len; i++ {
		if v.AtUnsafe(vctx, i) != o.AtUnsafe(octx, i) {
			return false
		}
	}

	return true
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

// Find the first item in the view bounds that returns true on the provided predicate.
// Return the index of such item (relative to the view start offset).
func (v UnmanagedView[T, Offset]) Index(ctx ViewContext[T], f func(T) bool) (Offset, error) {
	for idx, item := range v.Range2(ctx) {
		if f(item) {
			return idx, nil
		}
	}

	var o Offset
	return o, errors.New("no item found that matches predicate")
}

// Partition this view to two consecutive views, splitting them at the provided index.
func (v UnmanagedView[T, Offset]) Partition(
	ctx ViewContext[T], index Offset,
) (UnmanagedView[T, Offset], UnmanagedView[T, Offset]) {
	return v.Subview(0, index), v.Subview(index, v.Len())
}

// Similar to strings.FieldsFunc.
// Splits the input view at each run of items satisfying f(item) and returns an
// array of subviews of the origin view.
//
// Fields makes no guarantees about the order in which it calls f and assumes that
// f always outputs the same value for a given input.
func (v UnmanagedView[T, Offset]) Fields(ctx ViewContext[T], f func(T) bool) []UnmanagedView[T, Offset] {
	fields := make([]UnmanagedView[T, Offset], 0)
	start := Offset(0)
	collecting := false

	for end, item := range v.Range2(ctx) {
		shouldSplit := f(item)
		if shouldSplit && collecting {
			collecting = false
			subview := UnmanagedView[T, Offset]{Start: start, End: end}
			fields = append(fields, subview)
		} else if !shouldSplit && !collecting {
			collecting = true
			start = end
		}
	}

	if collecting {
		subview := UnmanagedView[T, Offset]{Start: start, End: v.Len()}
		fields = append(fields, subview)
	}

	return fields
}
