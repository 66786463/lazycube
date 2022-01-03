package length_test

import (
	"testing"

	"github.com/chucklebot/lazycube/length"
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
