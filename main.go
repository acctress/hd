package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var file string
var bpr int
var offsetlen int

func init() {
	flag.StringVar(&file, "file", "test.txt", "File")
	flag.IntVar(&bpr, "bpr", 16, "Bytes per row")
	flag.IntVar(&offsetlen, "ofl", 8, "Offset length")
}

func main() {
	flag.Parse()

	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("[hd:error] %v", err)
	}

	for i := 0; i < len(bytes); i += bpr {
		last := min(i+bpr, len(bytes))
		row := bytes[i:last]

		var hex_str strings.Builder
		for j := range len(row) {
			fmt.Fprintf(&hex_str, "%02x ", row[j])
		}

		padding := bpr - len(row)
		if padding > 0 {
			for k := range padding {
				_ = k
				fmt.Fprintf(&hex_str, ".. ")
			}
		}

		var ascii_str strings.Builder
		for j := range len(row) {
			if row[j] >= 32 && row[j] <= 126 {
				fmt.Fprintf(&ascii_str, "%c", row[j])
			} else {
				fmt.Fprint(&ascii_str, "..")
			}
		}

		fmt.Printf("\x1b[38;5;189m%0*d    \x1b[38;5;217m%s    \x1b[38;5;69m|%s|\x1b[0m\n", offsetlen, i, hex_str.String(), ascii_str.String())
	}
}
