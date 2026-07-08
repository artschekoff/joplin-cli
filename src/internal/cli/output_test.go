package cli

import "testing"

func TestMsToISO(t *testing.T) {
	if got := msToISO(0); got != nil {
		t.Fatalf("zero should be nil, got %v", *got)
	}
	got := msToISO(1704067200000) // 2024-01-01T00:00:00Z
	if got == nil || *got != "2024-01-01T00:00:00Z" {
		t.Fatalf("iso conversion wrong: %v", got)
	}
}
