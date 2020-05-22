package tag

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"git.maxset.io/web/knaxim/pkg/srverror"
)

func isChar(b byte) bool {
	if b >= 'a' && b <= 'z' {
		return true
	}
	if b >= 'A' && b <= 'Z' {
		return true
	}
	if b >= '0' && b <= '9' {
		return true
	}
	return false
}

// ScanWords causes a scanner to extract each alpha-numeric sequence
func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := -1
	for i, b := range data {
		if start < 0 && isChar(b) {
			start = i
			continue
		}
		if start >= 0 && !isChar(b) {
			return i, bytes.ToLower(data[start:i]), nil
		}
	}
	if start < 0 {
		return len(data), nil, nil
	}
	if atEOF {
		return len(data), bytes.ToLower(data[start:]), nil
	}
	return start, nil, nil
}

// ExtractContentTags generates an array of tags for each unique word as defined by ScanWords
func ExtractContentTags(content io.Reader) ([]Tag, error) {
	cache := make(map[string]Tag)

	sc := bufio.NewScanner(content)
	sc.Split(ScanWords)

	for sc.Scan() {
		w := sc.Text()
		if _, present := cache[w]; !present {
			cache[w] = Tag{
				Word: w,
				Type: CONTENT,
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, srverror.New(err, 500, "Error 501", "ExtractContentTags scanning")
	}

	out := make([]Tag, 0, len(cache))
	for _, v := range cache {
		out = append(out, v)
	}
	return out, nil
}

// BuildNameTags converts a string into tags of that string and substrings that are alpha numeric sequences
func BuildNameTags(s string) (out []Tag, err error) {
	out = append(out, Tag{
		Word: s,
		Type: NAME,
	})

	sc := bufio.NewScanner(strings.NewReader(s))
	sc.Split(ScanWords)
	wordcache := make(map[string]bool)
	for sc.Scan() {
		w := sc.Text()
		if !wordcache[w] {
			out = append(out, Tag{
				Word: w,
				Type: NAME,
			})
			wordcache[w] = true
		}
	}
	err = sc.Err()
	return
}
