package domain

import "testing"

func TestNewSankhyaLineagePreservesExplicitStates(t *testing.T) {
	origin := InternalDocumentLine{DocumentID: 31301, LineNumber: 2}
	quantity := func(value float64) *float64 { return &value }

	tests := []struct {
		name        string
		expected    *float64
		descendants []SankhyaDescendant
		want        SankhyaLineageState
		wantErr     bool
	}{
		{name: "none", descendants: nil, want: SankhyaLineageNone},
		{name: "complete one to many", expected: quantity(3), descendants: []SankhyaDescendant{
			{Identity: InternalDocumentLine{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: quantity(1)},
			{Identity: InternalDocumentLine{DocumentID: 30602, LineNumber: 4}, AttendedQuantity: quantity(2)},
		}, want: SankhyaLineageComplete},
		{name: "unknown expected quantity stays partial", descendants: []SankhyaDescendant{
			{Identity: InternalDocumentLine{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: quantity(1)},
		}, want: SankhyaLineagePartial},
		{name: "partial nullable quantity", descendants: []SankhyaDescendant{
			{Identity: InternalDocumentLine{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: nil},
		}, want: SankhyaLineagePartial},
		{name: "partial attended total", expected: quantity(4), descendants: []SankhyaDescendant{
			{Identity: InternalDocumentLine{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: quantity(2)},
		}, want: SankhyaLineagePartial},
		{name: "conflicting duplicate", descendants: []SankhyaDescendant{
			{Identity: InternalDocumentLine{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: quantity(1)},
			{Identity: InternalDocumentLine{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: quantity(2)},
		}, want: SankhyaLineageConflict, wantErr: true},
		{name: "identical duplicate", descendants: []SankhyaDescendant{
			{Identity: InternalDocumentLine{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: quantity(1)},
			{Identity: InternalDocumentLine{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: quantity(1)},
		}, want: SankhyaLineageConflict, wantErr: true},
		{name: "over attended", expected: quantity(1), descendants: []SankhyaDescendant{
			{Identity: InternalDocumentLine{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: quantity(2)},
		}, want: SankhyaLineageConflict, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewSankhyaLineage(origin, test.expected, test.descendants)
			if (err != nil) != test.wantErr {
				t.Fatalf("NewSankhyaLineage() error = %v, wantErr %v", err, test.wantErr)
			}
			if got.State != test.want {
				t.Fatalf("NewSankhyaLineage() state = %q, want %q", got.State, test.want)
			}
			if len(got.Descendants) != len(test.descendants) {
				t.Fatalf("NewSankhyaLineage() descendants = %d, want %d", len(got.Descendants), len(test.descendants))
			}
		})
	}
}

func TestNewSankhyaLineageRejectsZeroIdentityWithoutDefaulting(t *testing.T) {
	got, err := NewSankhyaLineage(InternalDocumentLine{}, nil, nil)
	if !IsReadErrorCode(err, ReadErrorLineageConflict) {
		t.Fatalf("error = %v, want lineage_conflict", err)
	}
	if got.Origin.DocumentID != 0 || got.State != "" {
		t.Fatalf("result = %+v, want unresolved invalid identity", got)
	}
}
