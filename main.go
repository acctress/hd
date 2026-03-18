package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var file string
var bpr int
var offsetlen int
var search string

func init() {
	flag.StringVar(&file, "file", "test.txt", "File")
	flag.StringVar(&search, "search", "", "Search for a pattern")
	flag.IntVar(&bpr, "bpr", 16, "Bytes per row")
	flag.IntVar(&offsetlen, "ofl", 8, "Offset length")
}

func main() {
	flag.Parse()

	file_bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("[hd:error] %v", err)
	}

	matches := [][2]int{}
	if search != "" {
		pos := 0
		for true {
			idx := bytes.Index(file_bytes[pos:], []byte(search))
			if idx == -1 {
				break
			}

			matches = append(matches, [2]int{pos + idx, pos + idx + len(search)})
			pos += idx + len(search)
		}
	}

	for i := 0; i < len(file_bytes); i += bpr {
		last := min(i+bpr, len(file_bytes))
		row := file_bytes[i:last]

		var hex_str strings.Builder
		for j := range len(row) {
			if search != "" {
				matched := false
				for _, match := range matches {
					if i+j >= match[0] && i+j < match[1] {
						fmt.Fprintf(&hex_str, "\x1b[38;5;227m%02x\x1b[0m ", row[j])
						matched = true
						break
					}
				}

				if !matched {
					fmt.Fprintf(&hex_str, "\x1b[38;5;217m%02x\x1b[0m ", row[j])
				}
			} else {
				fmt.Fprintf(&hex_str, "\x1b[38;5;217m%02x\x1b[0m ", row[j])
			}
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
				fmt.Fprint(&ascii_str, ".")
			}
		}

		fmt.Printf("\x1b[38;5;189m%0*x    %s    \x1b[38;5;69m|%-*s|\x1b[0m\n", offsetlen, i, hex_str.String(), bpr, ascii_str.String())
	}
}
