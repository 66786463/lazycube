// alen shows the lengths of the supplied files.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/66786463/lazycube/length"
)

var (
	doAccumulate = flag.Bool("accumulate", false, "show running total")
	doCheck      = flag.Bool("check", false, "check round sectors")
	doTotal      = flag.Bool("total", false, "show total length")
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	if *doCheck && *doAccumulate {
		log.Print("ignoring -accumulate since -check is set")
		*doAccumulate = false
	}
	if *doCheck && *doTotal {
		log.Print("ignoring -total since -check is set")
		*doTotal = false
	}
	total := &length.RawLength{}
	for _, f := range flag.Args() {
		rl, err := length.FetchLength(f)
		if err != nil {
			log.Fatalf("while fetching length of %s: %v", f, err)
		}
		if total.Rate != 0 && total.Rate != rl.Rate {
			msg := fmt.Sprintf("sample rate changed from %d to %d while processing %s",
				total.Rate, rl.Rate, f)
			if *doAccumulate || *doTotal {
				log.Fatalf(msg)
			} else {
				log.Print(msg)
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
