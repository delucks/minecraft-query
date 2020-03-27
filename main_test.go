package main

import (
	"fmt"
	"testing"
)

func TestFindByteSequence(t *testing.T) {
	idx := FindByteSequence([]byte{0x00, 0x00}, []byte{0x01, 0x02, 0x03, 0x00, 0x00, 0x09})
	if idx != 3 {
		t.Error(fmt.Sprintf("Double-0 not found in byte sequence. Returned %d", idx))
	}
}
