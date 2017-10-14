package script

import (
	"io"
	"os"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"gopkg.in/cheggaaa/pb.v1"
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

// ProgressReader returns a reader that is able to visualize read progress.
func (c Context) ProgressReader(reader io.Reader, size int) (io.Reader, *pb.ProgressBar) {
	bar := pb.New(size).SetUnits(pb.U_BYTES)

	// create proxy reader
	return bar.NewProxyReader(reader), bar
}

// ProgressFileReader returns a reader that is able to visualize read progress for a file.
func (c Context) ProgressFileReader(f *os.File) (io.Reader, *pb.ProgressBar, error) {
	size, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}
	r, bar := c.ProgressReader(f, int(size.Size()))
	return r, bar, nil
}
