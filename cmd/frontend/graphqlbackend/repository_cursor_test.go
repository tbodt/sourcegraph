package graphqlbackend

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	rawCursor    = repositoryCursor{After: 1}
	opaqueCursor = "UmVwb3NpdG9yeUN1cnNvcjp7IkFmdGVyIjoxfQ=="
)

func TestMarshalRepositoryCursor(t *testing.T) {
	if got, want := marshalRepositoryCursor(&rawCursor), opaqueCursor; got != want {
		t.Errorf("got opaque cursor %q, want %q", got, want)
	}
}

func TestUnmarshalRepositoryCursor(t *testing.T) {
	cursor, err := unmarshalRepositoryCursor(&opaqueCursor)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(cursor, &rawCursor); diff != "" {
		t.Fatal(diff)
	}
}
