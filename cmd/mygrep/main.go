package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2)
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}
	ok, err := matchLine(string(line), pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

}

func matchLine(line string, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}
	for j := 0; j < len(line); j++ {
		i := 0
		for  j < len(line) && i < len(pattern) {
			fmt.Println(string(pattern[i]), string(line[j]))
			if pattern[i] == '\\' {
				if pattern[i+1] == 'd' {
					if(!strings.Contains("1234567890", string(line[j]))) {
						i = len(pattern) + 2
						break
					}
					j++; i+=2
				} else if pattern[i+1] == 'w' {
					if !strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", string(line[j])) {
						i = len(pattern) + 2
						break
					}
					j++; i+=2
				}
			} else if pattern[i] == '[' {
				k := i
				for k < len(pattern) && pattern[k] != ']' {
					k++
				}
				if k == len(pattern) {
					return false , fmt.Errorf("invalid pattern: %q", pattern)
				}
				if pattern[i+1] == '^' {
					if strings.Contains(pattern[j:k], string(line[j])) {
						i = len(pattern) + 2
						break
					}
				} else {
					if !strings.Contains(pattern[1:len(pattern)-1], string(line[j])) {
						i = len(pattern) + 2
						break
					}
				}
				j = k + 1; i++
			} else {
				if pattern[i] != line[j] {
					i = len(pattern) + 2
					break
				}
				i++; j++
			}
		}

		if i == len(pattern) {
			return true, nil
		}
	}

	return false, nil
}
