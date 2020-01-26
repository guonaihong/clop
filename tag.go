package clop

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrNotPointerType = errors.New("Not pointer type")
var ErrUnsupportedType = errors.New("Unsupported type")

func contains(s string, b byte) bool {
	return strings.IndexByte(s, b) != -1
}

// 本函数来自stdlib, 简单修改了下
func Unquote(s string) (string, error) {
	n := len(s)
	if n < 2 {
		return "", strconv.ErrSyntax
	}
	quote := s[0]
	if quote != s[n-1] {
		return "", strconv.ErrSyntax
	}
	s = s[1 : n-1]

	if quote == '`' {
		if contains(s, '`') {
			return "", strconv.ErrSyntax
		}
		if contains(s, '\r') {
			// -1 because we know there is at least one \r to remove.
			buf := make([]byte, 0, len(s)-1)
			for i := 0; i < len(s); i++ {
				if s[i] != '\r' {
					buf = append(buf, s[i])
				}
			}
			return string(buf), nil
		}
		return s, nil
	}
	if quote != '"' && quote != '\'' {
		return "", strconv.ErrSyntax
	}
	/*
		if contains(s, '\n') {
			return "", strconv.ErrSyntax
		}
	*/

	// Is it trivial? Avoid allocation.
	if !contains(s, '\\') && !contains(s, quote) {
		switch quote {
		case '"':
			if utf8.ValidString(s) {
				return s, nil
			}
		case '\'':
			r, size := utf8.DecodeRuneInString(s)
			if size == len(s) && (r != utf8.RuneError || size != 1) {
				return s, nil
			}
		}
	}

	var runeTmp [utf8.UTFMax]byte
	buf := make([]byte, 0, 3*len(s)/2) // Try to avoid more allocations.
	for len(s) > 0 {
		c, multibyte, ss, err := strconv.UnquoteChar(s, quote)
		if err != nil {
			return "", err
		}
		s = ss
		if c < utf8.RuneSelf || !multibyte {
			buf = append(buf, byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf = append(buf, runeTmp[:n]...)
		}
		if quote == '\'' && len(s) != 0 {
			// single-quoted must be single character
			return "", strconv.ErrSyntax
		}
	}
	return string(buf), nil
}

type Tag string

// 该函数原版本来自stdlib, 简单修改了下
func (tag Tag) Lookup(key string) (value string, ok bool) {

	for tag != "" {
		i := 0
		for i < len(tag) && (tag[i] == ' ' || tag[i] == '\n' || tag[i] == '\t') {
			i++
		}

		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := Unquote(qvalue)
			//fmt.Printf("key(%s) name=(%s) tag(%s) qvalue(%s) err:%s\n", key, name, tag, qvalue, err)
			if err != nil {
				break
			}
			return value, true
		}
	}

	return "", false
}

func (tag Tag) Get(key string) (value string) {
	v, _ := tag.Lookup(key)
	return v
}
