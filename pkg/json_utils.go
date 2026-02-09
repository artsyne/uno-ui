package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalASCII marshals a value to JSON, escaping non-ASCII characters to \uXXXX.
// This is useful for putting JSON in HTTP headers where non-ASCII characters are not allowed.
func MarshalASCII(v interface{}) ([]byte, error) {
	// First marshal to standard JSON (UTF-8)
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Iterate and escape non-ASCII
	var buf bytes.Buffer
	// Pre-allocate assuming mostly ASCII
	buf.Grow(len(b))

	// We treat the byte slice as UTF-8 string
	// But simply iterating bytes is not enough, we need to decode runes
	runes := []rune(string(b))

	for _, r := range runes {
		if r < 128 {
			buf.WriteByte(byte(r))
		} else {
			if r > 0xFFFF {
				// Encode as surrogate pair for JSON compatibility
				r -= 0x10000
				s1 := 0xD800 + (r >> 10)
				s2 := 0xDC00 + (r & 0x3FF)
				fmt.Fprintf(&buf, "\\u%04x\\u%04x", s1, s2)
			} else {
				// Escape to \uXXXX
				fmt.Fprintf(&buf, "\\u%04x", r)
			}
		}
	}

	return buf.Bytes(), nil
}
