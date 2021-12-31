// alen shows the lengths of the supplied files.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jfreymuth/oggvorbis"
	"github.com/mewkiz/flac"
)

const (
	CDDARate          = 44100
	CDDASectorSamples = CDDARate / 75
)

type RawLength struct {
	Rate    uint32
	Samples uint64
}

func (rl *RawLength) String() string {
	return fmt.Sprintf("%d samples @ %d Hz", rl.Samples, rl.Rate)
}

func (rl *RawLength) ToCDDALength() *CDDALength {
	if rl.Samples == 0 {
		return &CDDALength{}
	}
	minutes := rl.Samples / (uint64(rl.Rate) * 60)
	cl := &CDDALength{
		Rate:    rl.Rate,
		Minutes: uint32(minutes),
		Seconds: uint32((rl.Samples - (uint64(minutes) * uint64(rl.Rate) * 60)) / uint64(rl.Rate)),
	}
	remainder := rl.Samples -
		((uint64(cl.Minutes) * uint64(cl.Rate) * 60) + (uint64(cl.Seconds) * uint64(cl.Rate)))
	if cl.Rate == CDDARate {
		cl.Sectors = uint32(remainder / CDDASectorSamples)
		cl.Samples = uint32(remainder % CDDASectorSamples)
	} else {
		cl.Samples = uint32(remainder)
	}
	return cl
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
	if cl.Rate == CDDARate {
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
	f, err := flac.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return &RawLength{
		Rate:    f.Info.SampleRate,
		Samples: f.Info.NSamples,
	}, nil
}

func fetchOggVorbisLength(path string) (*RawLength, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	samples, format, err := oggvorbis.GetLength(f)
	if err != nil {
		return nil, err
	}
	return &RawLength{
		Rate:    uint32(format.SampleRate),
		Samples: uint64(samples),
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
	for _, f := range flag.Args() {
		switch {
		case strings.HasSuffix(f, ".flac"):
			rl, err := fetchFLACLength(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				break
			}
			fmt.Printf("%s\t%s\n", rl.ToCDDALength(), f)
		case strings.HasSuffix(f, ".ogg"):
			rl, err := fetchOggVorbisLength(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				break
			}
			fmt.Printf("%s\t%s\n", rl.ToCDDALength(), f)
		default:
			fmt.Fprintf(os.Stderr, "W: we don't do whatever the hell %q is!\n", f)
		}
	}
}
