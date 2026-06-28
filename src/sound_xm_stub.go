package main

import (
	"fmt"
	"io"

	"github.com/gopxl/beep/v2"
)

func xmpDecode(f io.ReadSeekCloser) (beep.StreamSeekCloser, beep.Format, error) {
	return nil, beep.Format{}, fmt.Errorf("libxmp support is not built in; rebuild with -tags libxmp and libxmp development headers installed")
}
