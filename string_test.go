package view_test

import (
	"testing"

	"github.com/RealA10N/view"
	"github.com/stretchr/testify/assert"
)

func TestStringDetach(t *testing.T) {
	v := view.NewStringView[uint]("hello")
	unmanaged, ctx := v.Subview(1, 4).Detach()
	assert.Equal(t, "ell", unmanaged.String(ctx))
}
