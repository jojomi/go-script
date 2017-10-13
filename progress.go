package script

import (
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

// ActivityIndicatorCustom returns an activity indicator as specified.
// Use Start() to start and Stop() to stop it.
// See: https://godoc.org/github.com/gernest/wow
func (c Context) ActivityIndicatorCustom(text string, typ spin.Name) *wow.Wow {
	w := wow.New(c.stdout, spin.Get(typ), " "+text)
	return w
}

// ActivityIndicator returns an activity indicator with specified text and default animation.
func (c Context) ActivityIndicator(text string) *wow.Wow {
	return c.ActivityIndicatorCustom(text, spin.Dots)
}
