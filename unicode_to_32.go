package main

import (
	"fmt"
	"unicode/utf8"
)

func main3() {
	input := "hello, 世界"
	
	// Constrain input to 32 bytes (8 characters)
	limit := 32
	if utf8.RuneCountInString(input) > limit {
		input = truncateUTF8(input, limit)
	}
	fmt.Printf("Input constrained to %d bytes: %q\n", limit, input)
	
	// Constrain input to 64 bytes (16 characters)
	limit = 64
	if utf8.RuneCountInString(input) > limit {
		input = truncateUTF8(input, limit)
	}
	fmt.Printf("Input constrained to %d bytes: %q\n", limit, input)
}

// truncateUTF8 truncates a UTF-8 encoded string to the given byte length
func truncateUTF8(s string, maxBytes int) string {
	for len(s) > maxBytes {
		_, size := utf8.DecodeLastRuneInString(s)
		s = s[:len(s)-size]
	}
	return s
}
