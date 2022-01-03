package length

import (
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

func FetchLength(path string) (*RawLength, error) {
	switch {
	case strings.HasSuffix(path, ".flac"):
		return fetchFLACLength(path)
	case strings.HasSuffix(path, ".ogg"):
		return fetchOggVorbisLength(path)
	default:
		return nil, fmt.Errorf("don't know how to fetch length of %v", path)
	}
}
