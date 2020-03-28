package main

import (
	"fmt"
	"testing"
)

func TestFindByteSequence(t *testing.T) {
	idx := FindByteSequence([]byte{0x00, 0x00, 0x01, 0xde}, []byte{0x01})
	if idx != -1 {
		t.Error(fmt.Sprintf("A longer needle than haystack found an index (%d)", idx))
	}
	idx = FindByteSequence([]byte{0x00, 0x00}, []byte{0x01, 0x02, 0x03, 0x00, 0x00, 0x09})
	if idx != 3 {
		t.Error(fmt.Sprintf("Double-0 not found in byte sequence. Returned %d", idx))
	}
	idx = FindByteSequence([]byte{0x00, 0x00}, []byte{0x01, 0x02, 0x03, 0x09})
	if idx != -1 {
		t.Error(fmt.Sprintf("Found a nonexistent pattern in the byte sequence. Returned %d", idx))
	}
	idx = FindByteSequence([]byte{0x01}, []byte{0x01, 0x02, 0x03, 0x09})
	if idx != 0 {
		t.Error("Couldn't find the first byte in a sequence properly")
	}
	idx = FindByteSequence([]byte{0x02, 0x03}, []byte{0x01, 0x02, 0x03, 0x09})
	if idx != 1 {
		t.Error("Couldn't find middle sequence properly")
	}
}
