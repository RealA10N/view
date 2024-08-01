package view

import "golang.org/x/exp/constraints"

type UnmanagedStringView[Offset constraints.Unsigned] struct{ UnmanagedView[rune, Offset] }

func NewUnmanagedStringView[Offset constraints.Unsigned](s string) (UnmanagedStringView[Offset], ViewContext[rune]) {
	unmanaged, ctx := NewUnmanagedView[rune, Offset]([]rune(s))
	return UnmanagedStringView[Offset]{unmanaged}, ctx
}

func (v UnmanagedStringView[Offset]) String(ctx ViewContext[rune]) string {
	return string(v.Raw(ctx))
}

func (v UnmanagedStringView[Offset]) Attach(ctx ViewContext[rune]) StringView[Offset] {
	return StringView[Offset]{v.UnmanagedView.Attach(ctx)}
}

func (v UnmanagedStringView[Offset]) Subview(start, end Offset) UnmanagedStringView[Offset] {
	return UnmanagedStringView[Offset]{v.UnmanagedView.Subview(start, end)}
}

type StringView[Offset constraints.Unsigned] struct{ View[rune, Offset] }

func NewStringView[Offset constraints.Unsigned](s string) StringView[Offset] {
	view := NewView[rune, Offset]([]rune(s))
	return StringView[Offset]{view}
}

func (v StringView[T]) String() string {
	return string(v.unmanaged.Raw(v.ctx))
}

func (v StringView[Offset]) Detach() (UnmanagedStringView[Offset], ViewContext[rune]) {
	view, ctx := v.View.Detach()
	return UnmanagedStringView[Offset]{view}, ctx
}

func (v StringView[Offset]) Subview(start, end Offset) StringView[Offset] {
	return StringView[Offset]{v.View.Subview(start, end)}
}
