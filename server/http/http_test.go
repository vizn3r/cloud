package http_test

import (
	"testing"

	"cloud-server/http"
)

func TestIsTest(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := http.IsTest()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("IsTest() = %v, want %v", got, tt.want)
			}
		})
	}
}
