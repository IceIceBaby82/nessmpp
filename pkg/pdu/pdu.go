package pdu

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrCStringTooLong   = errors.New("c-string is too long")
	ErrInvalidCString   = errors.New("invalid c-string: contains null byte before end")
	ErrMissingNullByte  = errors.New("c-string is not null terminated")
	ErrInvalidMaxLength = errors.New("invalid max length: must be greater than 0")
)

// CString represents a null-terminated string used in SMPP PDUs
type CString string

// EncodeCString encodes a string as a C-string (null-terminated) with maximum length check
func EncodeCString(s string, maxLen int) ([]byte, error) {
	if maxLen <= 0 {
		return nil, ErrInvalidMaxLength
	}

	// Check if the string length (including null terminator) exceeds maxLen
	if len(s)+1 > maxLen {
		return nil, ErrCStringTooLong
	}

	// Check for null bytes within the string
	if bytes.IndexByte([]byte(s), 0) != -1 {
		return nil, ErrInvalidCString
	}

	// Create buffer with string + null terminator
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = 0 // Add null terminator

	return buf, nil
}

// DecodeCString decodes a C-string from a reader with maximum length limit
func DecodeCString(r io.Reader, maxLen int) (string, error) {
	if maxLen <= 0 {
		return "", ErrInvalidMaxLength
	}

	buf := make([]byte, maxLen)
	var result []byte

	for i := 0; i < maxLen; i++ {
		if _, err := io.ReadFull(r, buf[:1]); err != nil {
			if err == io.EOF {
				return "", ErrMissingNullByte
			}
			return "", err
		}

		if buf[0] == 0 {
			return string(result), nil
		}

		result = append(result, buf[0])
	}

	return "", ErrCStringTooLong
}

// String returns the string representation of the CString
func (cs CString) String() string {
	return string(cs)
}

// Bytes returns the byte slice representation of the CString including null terminator
func (cs CString) Bytes() []byte {
	b := make([]byte, len(cs)+1)
	copy(b, cs)
	b[len(cs)] = 0
	return b
}

// Length returns the length of the CString including the null terminator
func (cs CString) Length() int {
	return len(cs) + 1
}

// ReadCString reads a C-string from a byte slice
func ReadCString(data []byte) (string, int, error) {
	if len(data) == 0 {
		return "", 0, ErrMissingNullByte
	}

	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return "", 0, ErrMissingNullByte
	}

	return string(data[:nullIndex]), nullIndex + 1, nil
}

// WriteCString writes a C-string to a writer
func WriteCString(w io.Writer, s string) error {
	b := make([]byte, len(s)+1)
	copy(b, s)
	b[len(s)] = 0
	_, err := w.Write(b)
	return err
}
