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
func charMatches(pattern string, lineChar byte) bool {
	switch pattern {
	case "\\d":
		return strings.Contains("0123456789", string(lineChar))
	case "\\w":
		return strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", string(lineChar))
	default:
		return pattern[0] == lineChar
	}
}

func matchLine(line string, pattern string) (bool, error) {
    if utf8.RuneCountInString(pattern) == 0 {
        return false, fmt.Errorf("unsupported pattern: %q", pattern)
    }
    endAnchor := false
    length := len(line)
    if pattern[0] == '^' {
        length = 1
        pattern = pattern[1:]
    }
    if pattern[len(pattern)-1] == '$' {
        pattern = pattern[:len(pattern)-1]
        endAnchor = true
    }
    for j := 0; j < length; j++ {
        i := 0
        matchStart := j
        for j < len(line) && i < len(pattern) {
            if pattern[i] == '\\' {
                if i+1 < len(pattern) && pattern[i+1] == '\\' {
                    i++
                }
                if i+1 < len(pattern) {
                    if pattern[i+1] == 'd' || pattern[i+1] == 'w' {
                        if !charMatches(pattern[i:i+2], line[j]) {
                            break
                        }
                        j++
                        i += 2
                    } else {
                        return false, fmt.Errorf("invalid pattern: %q", pattern)
                    }
                }
            } else if pattern[i] == '[' {
				k := i
				for k < len(pattern) && pattern[k] != ']' {
					k++
				}
				if k == len(pattern) {
					return false, fmt.Errorf("invalid pattern: %q", pattern)
				}

				if pattern[i+1] == '^' {

					if strings.Contains(pattern[i+2:k], string(line[j])) {
						i = len(pattern) + 2
						break
					}
				} else {
					if !strings.Contains(pattern[i+1:k], string(line[j])) {
						i = len(pattern) + 2
						break
					}
				}
				i = k + 1
				j++
            } else if i+1 < len(pattern) && (pattern[i+1] == '+' || pattern[i+1] == '*' || pattern[i+1] == '?') {
                operator := pattern[i+1]
                char := pattern[i]
                i += 2
                count := 0
                for j < len(line) && charMatches(string(char), line[j]) {
                    j++
                    count++
                }
                if operator == '+' && count == 0 {
                    break
                }
                if operator == '?' && count > 1 {
                    j -= count - 1
                }
            } else if pattern[i] == '.' {
				j++
				i++
			}else {
                if pattern[i] != line[j] {
                    break
                }
                i++
                j++
            }
        }
        if i == len(pattern) {
            if endAnchor {
                if j == len(line) {
                    return true, nil
                }
            } else {
                return true, nil
            }
        }
        j = matchStart
    }
    return false, nil
}
