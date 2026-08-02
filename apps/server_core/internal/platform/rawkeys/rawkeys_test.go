package rawkeys

import (
	"encoding/json"
	"testing"
)

type sampleDTO struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Nested struct {
		Mode string `json:"mode"`
	} `json:"shipping"`
	NoTag    string
	Excluded string `json:"-"`
}

func TestUndeclaredNamesOnlyKeysTheDTOOmits(t *testing.T) {
	raw := json.RawMessage(`{"id":"MLB1","title":"x","shipping":{"mode":"me2"},"price":10,"currency_id":"BRL"}`)

	got, err := Undeclared(raw, sampleDTO{}, nil)
	if err != nil {
		t.Fatalf("Undeclared() error = %v", err)
	}

	want := []string{"currency_id", "price"}
	if len(got) != len(want) {
		t.Fatalf("Undeclared() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Undeclared() = %v, want %v", got, want)
		}
	}
}

func TestUndeclaredRespectsIgnoreList(t *testing.T) {
	raw := json.RawMessage(`{"id":"MLB1","title":"x","price":10}`)

	got, err := Undeclared(raw, sampleDTO{}, []string{"price"})
	if err != nil {
		t.Fatalf("Undeclared() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Undeclared() = %v, want empty", got)
	}
}

func TestUndeclaredRejectsNonObjectPayload(t *testing.T) {
	if _, err := Undeclared(json.RawMessage(`[1,2]`), sampleDTO{}, nil); err == nil {
		t.Fatal("Undeclared() error = nil, want error for non-object payload")
	}
}
