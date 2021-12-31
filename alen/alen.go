// alen shows the lengths of the supplied files.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/mewkiz/flac"
)

const (
	CDDA_SECTOR_SAMPLES = 75
)

type RawLength struct {
	Samples uint64
	Rate    int32
}

type CDDALength struct {
	Rate    int32
	Minutes int32
	Seconds int32
	Sectors int32
	Samples int32
}

var (
	doAccumulate = flag.Bool("accumulate", false, "show running total")
	doCheck      = flag.Bool("check", false, "check round sectors")
	doTotal      = flag.Bool("total", false, "show total length")
)

func main() {
	flag.Parse()
	if *doCheck && *doAccumulate {
		fmt.Fprintln(os.Stderr, "W: ignoring -accumulate since -check is set")
		*doAccumulate = false
	}
	if *doCheck && *doTotal {
		fmt.Fprintln(os.Stderr, "W: ignoring -total since -check is set")
		*doTotal = false
	}
	fmt.Println("Hello, world!")
	for _, f := range flag.Args() {
		switch {
		case strings.HasSuffix(f, ".flac"):
			fmt.Println("we do flac!")
		case strings.HasSuffix(f, ".ogg"):
			fmt.Println("we do ogg!")
		default:
			fmt.Fprintf(os.Stderr, "W: we don't do whatever the hell %q is!\n", f)
		}
	}
}
