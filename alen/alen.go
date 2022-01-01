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
	CDDARate         = 44100
	SamplesPerSector = CDDARate / 75
)

type RawLength struct {
	Rate    uint64
	Samples uint64
}

func (rl *RawLength) String() string {
	return fmt.Sprintf("%d samples @ %d Hz", rl.Samples, rl.Rate)
}

func (rl *RawLength) CDDALength() *CDDALength {
	cl := &CDDALength{
		Rate: rl.Rate,
	}
	if rl.Samples == 0 {
		return cl
	}
	cl.Minutes = rl.Samples / (rl.Rate * 60)
	cl.Seconds = (rl.Samples - (cl.Minutes * cl.Rate * 60)) / rl.Rate
	remainder := rl.Samples - ((cl.Minutes * cl.Rate * 60) + (cl.Seconds * cl.Rate))
	if cl.Rate == CDDARate {
		cl.Sectors = remainder / SamplesPerSector
		cl.Samples = remainder % SamplesPerSector
	} else {
		cl.Samples = remainder
	}
	return cl
}

type CDDALength struct {
	Rate    uint64
	Minutes uint64
	Seconds uint64
	Sectors uint64
	Samples uint64
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
		return nil, fmt.Errorf("%s: %v", path, err)
	}
	defer f.Close()
	return &RawLength{
		Rate:    uint64(f.Info.SampleRate),
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
		Rate:    uint64(format.SampleRate),
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
	total := &RawLength{}
	var err error
	for _, f := range flag.Args() {
		rl := &RawLength{}
		switch {
		case strings.HasSuffix(f, ".flac"):
			rl, err = fetchFLACLength(f)
		case strings.HasSuffix(f, ".ogg"):
			rl, err = fetchOggVorbisLength(f)
		default:
			fmt.Fprintf(os.Stderr, "E: we don't do whatever the hell %q is!\n", f)
			os.Exit(1)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "E: %v\n", err)
			os.Exit(1)
		}
		if total.Rate != 0 && total.Rate != rl.Rate {
			msg := fmt.Sprintf("sample rate changed from %d to %d while processing %s",
				total.Rate, rl.Rate, f)
			if *doAccumulate || *doTotal {
				fmt.Fprintf(os.Stderr, "E: %s; exiting\n", msg)
				os.Exit(1)
			} else {
				fmt.Fprintf(os.Stderr, "W: %s\n", msg)
			}
		}
		total.Rate = rl.Rate
		total.Samples += rl.Samples
		if *doTotal {
			continue
		}
		if *doAccumulate {
			fmt.Printf("%s\t%s\n", total.CDDALength(), f)
		} else {
			cl := rl.CDDALength()
			if *doCheck && cl.Samples == 0 {
				continue
			}
			fmt.Printf("%s\t%s\n", cl, f)
		}
	}
	if *doTotal {
		fmt.Printf("%s\t%s\n", total.CDDALength(), "total")
	}
}
