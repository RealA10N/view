package view

import "golang.org/x/exp/constraints"

type UnmanagedStringView[Offset constraints.Unsigned] struct{ UnmanagedView[rune, Offset] }

func (v UnmanagedStringView[Offset]) String(ctx ViewContext[rune]) string {
	return string(v.Raw(ctx))
}

type StringView[Offset constraints.Unsigned] struct{ View[rune, Offset] }

func (v StringView[T]) String() string {
	return string(v.unmanaged.Raw(v.ctx))
}
