package codec_test

import (
	"testing"

	"github.com/LYH2263/go-outboxrelay/internal/codec"
)

func TestCloneBytesIndependent(t *testing.T) {
	in := []byte("abc")
	out := codec.CloneBytes(in)
	in[0] = 'Z'
	if out[0] == 'Z' {
		t.Fatal("alias")
	}
}

func TestCloneHeadersIndependent(t *testing.T) {
	h := map[string]string{"a": "1"}
	c := codec.CloneHeaders(h)
	h["a"] = "2"
	if c["a"] != "1" {
		t.Fatal("header alias")
	}
}
