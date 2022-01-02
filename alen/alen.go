// alen shows the lengths of the supplied files.
package main

import (
	"flag"
	"fmt"
	"os"

	"ondioline.org/alen/length"
)

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
	total := &length.RawLength{}
	for _, f := range flag.Args() {
		rl, err := length.FetchLength(f)
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
