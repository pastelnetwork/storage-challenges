package helper

import "testing"

func TestGenerateFakePastelMnID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Fatal(GenerateFakePastelMnID())
		})
	}
}
