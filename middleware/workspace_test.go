package middleware

import "testing"

func TestParseWorkspaceHeader(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		wantID     uint
		wantExists bool
		wantErr    bool
	}{
		{name: "missing", value: "", wantExists: false},
		{name: "blank", value: "   ", wantExists: false},
		{name: "valid", value: "42", wantID: 42, wantExists: true},
		{name: "valid with spaces", value: " 7 ", wantID: 7, wantExists: true},
		{name: "zero", value: "0", wantExists: true, wantErr: true},
		{name: "negative", value: "-1", wantExists: true, wantErr: true},
		{name: "malformed", value: "abc", wantExists: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotExists, err := parseWorkspaceHeader(tt.value)
			if gotID != tt.wantID {
				t.Fatalf("workspace id = %d, want %d", gotID, tt.wantID)
			}
			if gotExists != tt.wantExists {
				t.Fatalf("exists = %t, want %t", gotExists, tt.wantExists)
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}
