package sourcekind

import "testing"

func TestSourceKind_Valid(t *testing.T) {
	cases := []struct {
		name string
		kind SourceKind
		want bool
	}{
		{"UploadSnapshot", UploadSnapshot, true},
		{"LiveReadThrough", LiveReadThrough, true},
		{"empty", SourceKind(""), false},
		{"bogus", SourceKind("bogus"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.kind.Valid(); got != tc.want {
				t.Errorf("SourceKind(%q).Valid() = %v, want %v", string(tc.kind), got, tc.want)
			}
		})
	}
}

func TestSourceKind_StringValues(t *testing.T) {
	if UploadSnapshot != "upload_snapshot" {
		t.Errorf("UploadSnapshot = %q, want %q", string(UploadSnapshot), "upload_snapshot")
	}
	if LiveReadThrough != "live_read_through" {
		t.Errorf("LiveReadThrough = %q, want %q", string(LiveReadThrough), "live_read_through")
	}
}
