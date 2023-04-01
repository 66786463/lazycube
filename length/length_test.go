package length_test

import (
	"testing"

	"github.com/66786463/lazycube/length"
)

func TestCDDALength(t *testing.T) {
	table := []struct {
		have *length.RawLength
		want *length.CDDALength
	}{
		{
			have: &length.RawLength{
				Rate:    44100,
				Samples: 0,
			},
			want: &length.CDDALength{
				Rate: 44100,
			},
		},
		{
			have: &length.RawLength{
				Rate:    44100,
				Samples: 588,
			},
			want: &length.CDDALength{
				Rate:    44100,
				Sectors: 1,
			},
		},
		{
			have: &length.RawLength{
				Rate:    44100,
				Samples: 9700524,
			},
			want: &length.CDDALength{
				Rate:    44100,
				Minutes: 3,
				Seconds: 39,
				Sectors: 72,
				Samples: 288,
			},
		},
		{
			have: &length.RawLength{
				Rate:    48000,
				Samples: 0,
			},
			want: &length.CDDALength{
				Rate: 48000,
			},
		},
		{
			have: &length.RawLength{
				Rate:    48000,
				Samples: 588,
			},
			want: &length.CDDALength{
				Rate:    48000,
				Samples: 588,
			},
		},
		{
			have: &length.RawLength{
				Rate:    44100,
				Samples: 159509700,
			},
			want: &length.CDDALength{
				Rate:    44100,
				Minutes: 60,
				Seconds: 17,
			},
		},
		{
			have: &length.RawLength{
				Rate:    44100,
				Samples: 3686827,
			},
			want: &length.CDDALength{
				Rate:    44100,
				Minutes: 1,
				Seconds: 23,
				Sectors: 45,
				Samples: 67,
			},
		},
		{
			have: &length.RawLength{
				Rate:    48000,
				Samples: 13360056,
			},
			want: &length.CDDALength{
				Rate:    48000,
				Minutes: 4,
				Seconds: 38,
				Samples: 16056,
			},
		},
	}

	for _, test := range table {
		if got := test.have.CDDALength(); *got != *test.want {
			t.Fatalf(`got "%v", want "%v"`, got, test.want)
		}
	}
}

func TestFetchLength(t *testing.T) {
	table := []struct {
		have string
		want *length.RawLength
	}{
		{
			have: "testdata/sin.440Hz@44100.flac",
			want: &length.RawLength{
				Rate:    44100,
				Samples: 1234567,
			},
		},
		{
			have: "testdata/sin.440Hz@44100.ogg",
			want: &length.RawLength{
				Rate:    44100,
				Samples: 1234567,
			},
		},
		{
			have: "testdata/sin.440Hz@48000.flac",
			want: &length.RawLength{
				Rate:    48000,
				Samples: 1234567,
			},
		},
		{
			have: "testdata/sin.440Hz@48000.flac",
			want: &length.RawLength{
				Rate:    48000,
				Samples: 1234567,
			},
		},
	}

	for _, test := range table {
		got, err := length.FetchLength(test.have)
		switch {
		case err != nil:
			t.Fatalf(`want "%v", got error: %v`, test.want, err)
		case *got != *test.want:
			t.Fatalf(`want "%v", got "%v"`, test.want, got)
		}
	}
}
