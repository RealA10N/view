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
