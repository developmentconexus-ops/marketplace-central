package domain

import "testing"

func TestIsValidGTIN(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "valid GTIN-8", value: "96385074", want: true},
		{name: "valid UPC-A GTIN-12", value: "036000291452", want: true},
		{name: "valid EAN-13 GTIN-13", value: "4006381333931", want: true},
		{name: "valid GTIN-14", value: "10012345678902", want: true},
		{name: "valid with surrounding whitespace", value: " 4006381333931\n", want: true},
		{name: "wrong check digit", value: "4006381333932", want: false},
		{name: "invalid GTIN-12 checksum", value: "123456789011", want: false},
		{name: "empty", value: "", want: false},
		{name: "whitespace only", value: " \t\n", want: false},
		{name: "wrong length 7 digits", value: "9638507", want: false},
		{name: "wrong length 9 digits", value: "963850740", want: false},
		{name: "wrong length 11 digits", value: "03600029145", want: false},
		{name: "letters", value: "1234567890abc", want: false},
		{name: "internal space", value: "40063813339 1", want: false},
		{name: "hyphen requires rejection", value: "4006381-333931", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidGTIN(tt.value); got != tt.want {
				t.Errorf("IsValidGTIN(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
