// alen shows the lengths of the supplied files.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mewkiz/flac/meta"
)

const (
	CDDA_SECTOR_SAMPLES = 75
)

type RawLength struct {
	Rate    uint32
	Samples uint64
}

func (rl *RawLength) String() string {
	return fmt.Sprintf("%d samples @ %d Hz", rl.Samples, rl.Rate)
}

type CDDALength struct {
	Rate    uint32
	Minutes uint32
	Seconds uint32
	Sectors uint32
	Samples uint32
}

func (cl *CDDALength) String() string {
	s := fmt.Sprintf("%2d:%02d", cl.Minutes, cl.Seconds)
	if cl.Rate == 44100 {
		s = fmt.Sprintf("%s.%02d", s, cl.Sectors)
	} else {
		s += "   "
	}
	if cl.Samples > 0 {
		s = fmt.Sprintf("%s +%d", s, cl.Samples)
	}
	return s
}

var (
	doAccumulate = flag.Bool("accumulate", false, "show running total")
	doCheck      = flag.Bool("check", false, "check round sectors")
	doTotal      = flag.Bool("total", false, "show total length")
)

func fetchFLACLength(path string) (*RawLength, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := meta.New(f)
	if err != nil {
		return nil, err
	}
	if b.Header.Type != meta.TypeStreamInfo {
		// NOTE(pjdc): Some malformed files don't put STREAMINFO first.  Handle?
		return nil, fmt.Errorf("first block in %q is not STREAMINFO (type %q)", path, b.Header.Type)
	}
	s, ok := b.Body.(meta.StreamInfo)
	if !ok {
		panic(fmt.Sprintf("%q: failed to cast block body to STREAMINFO; b = %v", path, b))
	}
	return &RawLength{
		Rate:    s.SampleRate,
		Samples: s.NSamples,
	}, nil
}

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
			rl, err := fetchFLACLength(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				break
			}
			fmt.Printf("%q: %s\n", f, rl)
		case strings.HasSuffix(f, ".ogg"):
			fmt.Println("we do ogg!")
		default:
			fmt.Fprintf(os.Stderr, "W: we don't do whatever the hell %q is!\n", f)
		}
	}
}
