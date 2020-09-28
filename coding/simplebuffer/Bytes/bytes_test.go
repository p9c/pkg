package Bytes

import (
	"testing"
)

func TestBytes(t *testing.T) {
	by := []byte("this is a test")
	bt := New()
	bt.Put(by)
	bt2 := New()
	bt2.Decode(bt.Encode())
	if string(bt.Get()) != string(bt2.Get()) {
		t.Fail()
	}
}
