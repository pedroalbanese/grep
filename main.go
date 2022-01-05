package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"common"
)

var (
	helpFlag = flag.Bool("help", false, "Show this help")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 || *helpFlag {
		println("`grep` <pattern> [<file>...]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	pattern, err := regexp.Compile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	if flag.NArg() == 1 {
		doGrep(pattern, os.Stdin, "<stdin>", false)
	} else {
		for _, fn := range flag.Args()[1:] {
			if fh, err := os.Open(fn); err == nil {
				func() {
					defer fh.Close()
					doGrep(pattern, fh, fn, flag.NArg() > 2)
				}()
			} else {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", fn, err)
			}
		}
	}
}

func doGrep(pattern *regexp.Regexp, fh io.Reader, fn string, print_fn bool) {
	buf := common.NewBufferedReader(fh)

	for {
		line, err := buf.ReadWholeLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading from %s: %v\n", fn, err)
			return
		}
		if line == "" {
			continue
		}

		if pattern.MatchString(line) {
			if print_fn {
				fmt.Printf("%s:", fn)
			}
			fmt.Printf("%s\n", line)
		}
	}
}