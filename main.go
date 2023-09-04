package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

var (
	helpFlag = flag.Bool("h", false, "Show this help")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 || *helpFlag || flag.Arg(0) == "-h" {
		fmt.Println("`grep` <pattern> [<file>...]")
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

func doGrep(pattern *regexp.Regexp, fh io.Reader, fn string, printFn bool) {
	buf := make([]byte, 4096)

	for {
		n, err := fh.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading from %s: %v\n", fn, err)
			return
		}
		if n == 0 {
			continue
		}

		lines := splitLines(buf[:n])
		for _, line := range lines {
			if pattern.MatchString(line) {
				if printFn {
					fmt.Printf("%s:", fn)
				}
				fmt.Printf("%s\n", line)
			}
		}
	}
}

func splitLines(data []byte) []string {
	var lines []string
	start := 0

	for i, b := range data {
		if b == '\n' {
			lines = append(lines, string(data[start:i]))
			start = i + 1
		}
	}

	if start < len(data) {
		lines = append(lines, string(data[start:]))
	}

	return lines
}
