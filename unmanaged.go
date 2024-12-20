package view

import (
	"errors"
	"iter"

	"golang.org/x/exp/constraints"
)

// The most basic slice view type.
// Internally, it contains a pointer to a heap allocated slice of type []T,
// and Start & End indices of type Offset.
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

// Returns the first item in the view bounds.
// If the view is empty, an error is returned, with an undefined value.
func (v UnmanagedView[T, Offset]) Front(ctx ViewContext[T]) (T, error) {
	if v.End <= v.Start {
		var t T
		return t, errors.New("view is empty")
	}

	return ctx[v.Start], nil
}

// Returns the first item in the view bounds.
// If the view is empty, the function might panic or return an undefined value.
func (v UnmanagedView[T, Offset]) FrontUnsafe(ctx ViewContext[T]) T {
	return ctx[v.Start]
}

// Returns the last item in the view bounds.
// If the view is empty, an error is returned, with an undefined value.
func (v UnmanagedView[T, Offset]) Back(ctx ViewContext[T]) (T, error) {
	if v.End <= v.Start {
		var t T
		return t, errors.New("view is empty")
	}

	return ctx[v.End-1], nil
}

// Returns the last item in the view bounds.
// If the view is empty, the function might panic or return an undefined value.
func (v UnmanagedView[T, Offset]) BackUnsafe(ctx ViewContext[T]) T {
	return ctx[v.End-1]
}

// Return a subview of the current view.
// Start and End indices are relative to the current view bounds,
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
	vctx ViewContext[T], u UnmanagedView[T, Offset], uctx ViewContext[T],
) bool {
	if v.Len() != u.Len() {
		return false
	}

	len := v.Len()
	var i Offset = 0
	for ; i < len; i++ {
		if v.AtUnsafe(vctx, i) != u.AtUnsafe(uctx, i) {
			return false
		}
	}

	return true
}

// Iterate over all values in the view (rangefunc).
func (v UnmanagedView[T, Offset]) Range(ctx ViewContext[T]) iter.Seq[T] {
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
func (v UnmanagedView[T, Offset]) Range2(ctx ViewContext[T]) iter.Seq2[Offset, T] {
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

// Find the first item in the view bounds that equals to the provided item.
// Return the index of such item (relative to the view start offset).
//
// If no items return true on the provided predicate, returns v.Len().
func (v UnmanagedView[T, Offset]) Index(ctx ViewContext[T], item T) Offset {
	for idx, cur := range v.Range2(ctx) {
		if cur == item {
			return idx
		}
	}

	return v.Len()
}

// Find the first item in the view bounds that returns true on the provided predicate.
// Return the index of such item (relative to the view start offset).
//
// If no items return true on the provided predicate, returns v.Len().
func (v UnmanagedView[T, Offset]) IndexFunc(ctx ViewContext[T], f func(T) bool) Offset {
	for idx, item := range v.Range2(ctx) {
		if f(item) {
			return idx
		}
	}

	return v.Len()
}

// Returns true iff the view contains the provided item.
func (v UnmanagedView[T, Offset]) Contains(ctx ViewContext[T], item T) bool {
	for cur := range v.Range(ctx) {
		if cur == item {
			return true
		}
	}
	return false
}

// Returns true iff the provided view is a prefix of the current view.
func (v UnmanagedView[T, Offset]) HasPrefix(
	ctx ViewContext[T],
	prefix UnmanagedView[T, Offset],
	prefixCtx ViewContext[T],
) bool {
	if v.Len() < prefix.Len() {
		return false
	}

	n := prefix.Len()
	for i := Offset(0); i < n; i++ {
		if v.AtUnsafe(ctx, i) != prefix.AtUnsafe(prefixCtx, i) {
			return false
		}
	}

	return true
}

// Returns true iff the provided view is a suffix of the current view.
func (v UnmanagedView[T, Offset]) HasSuffix(
	ctx ViewContext[T],
	suffix UnmanagedView[T, Offset],
	suffixCtx ViewContext[T],
) bool {
	if v.Len() < suffix.Len() {
		return false
	}

	n := suffix.Len()
	for i := Offset(0); i < n; i++ {
		if v.AtUnsafe(ctx, v.Len()-n+i) != suffix.AtUnsafe(suffixCtx, i) {
			return false
		}
	}

	return true
}

// Returns the longest common prefix of the current view and the provided one.
func (v UnmanagedView[T, Offset]) LongestCommonPrefix(
	ctx ViewContext[T],
	u UnmanagedView[T, Offset],
	uctx ViewContext[T],
) UnmanagedView[T, Offset] {
	n := min(v.Len(), u.Len())
	for i := Offset(0); i < n; i++ {
		if v.AtUnsafe(ctx, i) != u.AtUnsafe(uctx, i) {
			return v.Subview(0, i)
		}
	}
	return v.Subview(0, n)
}

// Returns the longest common suffix of the current view and the provided one.
func (v UnmanagedView[T, Offset]) LongestCommonSuffix(
	ctx ViewContext[T],
	u UnmanagedView[T, Offset],
	uctx ViewContext[T],
) UnmanagedView[T, Offset] {
	n := min(v.Len(), u.Len())
	for i := Offset(0); i < n; i++ {
		if v.AtUnsafe(ctx, v.Len()-i-1) != u.AtUnsafe(uctx, u.Len()-i-1) {
			return v.Subview(v.Len()-i, v.Len())
		}
	}
	return v.Subview(v.Len()-n, v.Len())
}

// Merge this and the other provided view into a one bigger view.
// This is done by setting newView.Start to min(v.Start, o.Start) and
// newView.End to max(v.End, o.End).
func (v UnmanagedView[T, Offset]) Merge(others ...UnmanagedView[T, Offset]) UnmanagedView[T, Offset] {
	nv := v

	for _, o := range others {
		nv.Start = min(nv.Start, o.Start)
		nv.End = max(nv.End, o.End)
	}

	return nv
}

// Merge this and the other provided view into a one bigger view, by returning
// a new view with the same end location, but the minimal start location out of
// all provided views.
func (v UnmanagedView[T, Offset]) MergeStart(others ...UnmanagedView[T, Offset]) UnmanagedView[T, Offset] {
	nv := v

	for _, o := range others {
		nv.Start = min(nv.Start, o.Start)
	}

	return nv
}

// Merge this and the other provided view into a one bigger view, by returning
// a new view with the same start location, but the maximal end location out of
// all provided views.
func (v UnmanagedView[T, Offset]) MergeEnd(others ...UnmanagedView[T, Offset]) UnmanagedView[T, Offset] {
	nv := v

	for _, o := range others {
		nv.End = max(nv.End, o.End)
	}

	return nv
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
