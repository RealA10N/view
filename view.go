package view

import (
	"golang.org/x/exp/constraints"
)

type View[T comparable, Offset constraints.Unsigned] struct {
	unmanaged UnmanagedView[T, Offset]
	ctx       ViewContext[T]
}

// Create a new (managed) view from an already existing slice.
// The view initially spans over the whole slice.
func NewView[T comparable, Offset constraints.Unsigned](data []T) View[T, Offset] {
	unmanaged, ctx := NewUnmanagedView[T, Offset](data)
	return unmanaged.Attach(ctx)
}

// Extract the unmanged view and context from the current view, and return
// copies of them. The old view is still valid and safe to use.
func (v View[T, Offset]) Detach() (UnmanagedView[T, Offset], ViewContext[T]) {
	return UnmanagedView[T, Offset](v.unmanaged), v.ctx
}

// Returns the raw underlying slice that the view is bound to.
func (v View[T, Offset]) Raw() []T {
	return v.unmanaged.Raw(v.ctx)
}

// Returns the size of the view slice.
func (v View[T, Offset]) Len() Offset {
	return v.unmanaged.Len()
}

// Returns the item at the provided index, relative to the view bounds.
// If the provided index goes out of the view bounds, an error is returned,
// with an undefined value.
func (v View[T, Offset]) At(index Offset) (T, error) {
	return v.unmanaged.At(v.ctx, index)
}

// Returns the item at the provided index, relative to the view bounds.
// Does not check view bounds, and if the provided index is greater than the
// view length, the function might panic or return an undefined value.
func (v View[T, Offset]) AtUnsafe(index Offset) T {
	return v.unmanaged.AtUnsafe(v.ctx, index)
}

// Return a subview of the current view, by start and end indecies.
func (v View[T, Offset]) Subview(start, end Offset) View[T, Offset] {
	return v.unmanaged.Subview(start, end).Attach(v.ctx)
}

// returns true if the underlying views are identical in their content.
func (v View[T, Offset]) Equal(o View[T, Offset]) bool {
	return v.unmanaged.Equal(v.ctx, o.unmanaged, o.ctx)
}

// Find the first item in the view bounds that returns true on the provided predicate.
// Return the index of such item (relative to the view start offset).
//
// If no items return true on the provided predicate, returns v.Len().
func (v View[T, Offset]) Index(f func(T) bool) Offset {
	return v.unmanaged.Index(v.ctx, f)
}

// Merge this and the other provided view into a one bigger view.
// This is done by setting newView.Start to min(v.Start, o.Start) and
// newView.End to max(v.End, o.End).
//
// This assumes that both views operate under the same context.
// More specificly, the context of the returned view will be the context of
// this view.
func (v View[T, Offset]) Merge(o View[T, Offset]) View[T, Offset] {
	return v.unmanaged.Merge(o.unmanaged).Attach(v.ctx)
}

// Partition this view to two consecutive views, splitting them at the provided index.
func (v View[T, Offset]) Partition(index Offset) (View[T, Offset], View[T, Offset]) {
	a, b := v.unmanaged.Partition(v.ctx, index)
	return a.Attach(v.ctx), b.Attach(v.ctx)
}

// Iterate over all values in the view (rangefunc).
func (v View[T, Offset]) Range() func(func(T) bool) {
	return v.unmanaged.Range(v.ctx)
}

// Iterate over all values in the view (rangefunc).
// Additionally, provides the iteration index as the first yield argument,
// where the index is relative to the view start.
func (v View[T, Offset]) Range2() func(func(Offset, T) bool) {
	return v.unmanaged.Range2(v.ctx)
}

// Similar to strings.FieldsFunc.
// Splits the input view at each run of items satisfying f(item) and returns an
// array of subviews of the origin view.
//
// Fields makes no guarantees about the order in which it calls f and assumes that
// f always outputs the same value for a given input.
func (v View[T, Offset]) Fields(f func(T) bool) []View[T, Offset] {
	unmanagedFields := v.unmanaged.Fields(v.ctx, f)
	fields := make([]View[T, Offset], len(unmanagedFields))
	for idx, unmanaged := range unmanagedFields {
		fields[idx] = unmanaged.Attach(v.ctx)
	}
	return fields
}
