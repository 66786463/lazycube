package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCDDALength(t *testing.T) {
	table := []struct {
		have *RawLength
		want *CDDALength
	}{
		{
			have: &RawLength{
				Rate:    44100,
				Samples: 0,
			},
			want: &CDDALength{
				Rate: 44100,
			},
		},
		{
			have: &RawLength{
				Rate:    44100,
				Samples: 588,
			},
			want: &CDDALength{
				Rate:    44100,
				Sectors: 1,
			},
		},
		{
			have: &RawLength{
				Rate:    44100,
				Samples: 9700524,
			},
			want: &CDDALength{
				Rate:    44100,
				Minutes: 3,
				Seconds: 39,
				Sectors: 72,
				Samples: 288,
			},
		},
		{
			have: &RawLength{
				Rate:    48000,
				Samples: 0,
			},
			want: &CDDALength{
				Rate: 48000,
			},
		},
		{
			have: &RawLength{
				Rate:    48000,
				Samples: 588,
			},
			want: &CDDALength{
				Rate:    48000,
				Samples: 588,
			},
		},
		{
			have: &RawLength{
				Rate:    44100,
				Samples: 159509700,
			},
			want: &CDDALength{
				Rate:    44100,
				Minutes: 60,
				Seconds: 17,
			},
		},
		{
			have: &RawLength{
				Rate:    48000,
				Samples: 13360056,
			},
			want: &CDDALength{
				Rate:    48000,
				Minutes: 4,
				Seconds: 38,
				Samples: 16056,
			},
		},
	}

	for _, test := range table {
		if got := test.have.CDDALength(); !cmp.Equal(got, test.want) {
			t.Fatalf(`got "%v", want "%v"`, got, test.want)
		}
	}
}
