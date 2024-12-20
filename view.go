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

// Extract the unmanaged view and context from the current view, and return
// copies of them. The old view is still valid and safe to use.
func (v View[T, Offset]) Detach() (UnmanagedView[T, Offset], ViewContext[T]) {
	return UnmanagedView[T, Offset](v.unmanaged), v.ctx
}

// Return the corresponding unmanaged view of the current view.
func (v View[T, Offset]) Unmanaged() UnmanagedView[T, Offset] {
	return v.unmanaged
}

// Return the context of the view.
func (v View[T, Offset]) Ctx() ViewContext[T] {
	return v.ctx
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

// Returns the first item in the view bounds.
// If the view is empty, an error is returned, with an undefined value.
func (v View[T, Offset]) Front() (T, error) {
	return v.unmanaged.Front(v.ctx)
}

// Returns the first item in the view bounds.
// If the view is empty, the function might panic or return an undefined value.
func (v View[T, Offset]) FrontUnsafe() T {
	return v.unmanaged.FrontUnsafe(v.ctx)
}

// Returns the last item in the view bounds.
// If the view is empty, an error is returned, with an undefined value.
func (v View[T, Offset]) Back() (T, error) {
	return v.unmanaged.Back(v.ctx)
}

// Returns the last item in the view bounds.
// If the view is empty, the function might panic or return an undefined value.
func (v View[T, Offset]) BackUnsafe() T {
	return v.unmanaged.BackUnsafe(v.ctx)
}

// Return a subview of the current view, by start and end indices.
func (v View[T, Offset]) Subview(start, end Offset) View[T, Offset] {
	return v.unmanaged.Subview(start, end).Attach(v.ctx)
}

// returns true if the underlying views are identical in their content.
func (v View[T, Offset]) Equal(o View[T, Offset]) bool {
	return v.unmanaged.Equal(v.ctx, o.unmanaged, o.ctx)
}

// Find the first item in the view bounds that equals to the provided item.
// Return the index of such item (relative to the view start offset).
//
// If no items return true on the provided predicate, returns v.Len().
func (v View[T, Offset]) Index(item T) Offset {
	return v.unmanaged.Index(v.ctx, item)
}

// Find the first item in the view bounds that returns true on the provided predicate.
// Return the index of such item (relative to the view start offset).
//
// If no items return true on the provided predicate, returns v.Len().
func (v View[T, Offset]) IndexFunc(f func(T) bool) Offset {
	return v.unmanaged.IndexFunc(v.ctx, f)
}

// Returns true iff the view contains the provided item.
func (v View[T, Offset]) Contains(item T) bool {
	return v.unmanaged.Contains(v.ctx, item)
}

// Returns true iff the provided view is a prefix of the current view.
func (v View[T, Offset]) HasPrefix(prefix View[T, Offset]) bool {
	unmanagedPrefix, prefixCtx := prefix.Detach()
	return v.unmanaged.HasPrefix(v.ctx, unmanagedPrefix, prefixCtx)
}

// Returns true iff the provided view is a suffix of the current view.
func (v View[T, Offset]) HasSuffix(suffix View[T, Offset]) bool {
	unmanagedSuffix, suffixCtx := suffix.Detach()
	return v.unmanaged.HasSuffix(v.ctx, unmanagedSuffix, suffixCtx)
}

// Returns the longest common prefix of the current view and the provided view.
func (v View[T, Offset]) LongestCommonPrefix(u View[T, Offset]) View[T, Offset] {
	return v.unmanaged.LongestCommonPrefix(v.ctx, u.unmanaged, u.ctx).Attach(v.ctx)
}

// Returns the longest common suffix of the current view and the provided view.
func (v View[T, Offset]) LongestCommonSuffix(u View[T, Offset]) View[T, Offset] {
	return v.unmanaged.LongestCommonSuffix(v.ctx, u.unmanaged, u.ctx).Attach(v.ctx)
}

// Merge this and the other provided view into a one bigger view.
// This is done by setting newView.Start to min(v.Start, o.Start) and
// newView.End to max(v.End, o.End).
//
// This assumes that both views operate under the same context.
// More specifically, the context of the returned view will be the context of
// this view.
func (v View[T, Offset]) Merge(others ...View[T, Offset]) View[T, Offset] {
	return v.unmanaged.Merge(detachMany(others)...).Attach(v.ctx)
}

// Merge this and the other provided view into a one bigger view, by returning
// a new view with the same end location, but the minimal start location out of
// all provided views.
func (v View[T, Offset]) MergeStart(others ...View[T, Offset]) View[T, Offset] {
	return v.unmanaged.MergeStart(detachMany(others)...).Attach(v.ctx)
}

// Merge this and the other provided view into a one bigger view, by returning
// a new view with the same start location, but the maximal end location out of
// all provided views.
func (v View[T, Offset]) MergeEnd(others ...View[T, Offset]) View[T, Offset] {
	return v.unmanaged.MergeEnd(detachMany(others)...).Attach(v.ctx)
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
	return attachMany(v.ctx, v.unmanaged.Fields(v.ctx, f))
}

func attachMany[T comparable, Offset constraints.Unsigned](
	ctx ViewContext[T], many []UnmanagedView[T, Offset],
) []View[T, Offset] {
	views := make([]View[T, Offset], len(many))
	for idx, unmanaged := range many {
		views[idx] = unmanaged.Attach(ctx)
	}
	return views
}

func detachMany[T comparable, Offset constraints.Unsigned](
	many []View[T, Offset],
) []UnmanagedView[T, Offset] {
	views := make([]UnmanagedView[T, Offset], len(many))
	for idx, view := range many {
		views[idx] = view.unmanaged
	}
	return views
}
