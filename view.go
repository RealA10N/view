package view

import (
	"golang.org/x/exp/constraints"
)

type View[T comparable, Offset constraints.Unsigned] struct {
	unmanaged UnmanagedView[T, Offset]
	ctx       ViewContext[T]
}

func (v UnmanagedView[T, Offset]) Attach(ctx ViewContext[T]) View[T, Offset] {
	return View[T, Offset]{
		unmanaged: v,
		ctx:       ctx,
	}
}

func NewView[T comparable, Offset constraints.Unsigned](data []T) View[T, Offset] {
	unmanaged, ctx := NewUnmanagedView[T, Offset](data)
	return unmanaged.Attach(ctx)
}

func (v View[T, Offset]) Raw() []T {
	return v.unmanaged.Raw(v.ctx)
}

// Returns the size of the view slice.
func (v View[T, Offset]) Len() Offset {
	return v.unmanaged.Len()
}

// Returns the item at the provided index, relative to the view bounds.
func (v View[T, Offset]) At(index Offset) (T, error) {
	return v.unmanaged.At(v.ctx, index)
}

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
