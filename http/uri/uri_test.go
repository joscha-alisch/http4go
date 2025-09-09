package uri

import (
	"errors"
	"testing"
)
import "github.com/google/go-cmp/cmp"

func TestUriOf(t *testing.T) {
	tests := []struct {
		name          string
		in            string
		want          Uri
		expectedError error
	}{
		{"simple", "http://example.com", Uri{
			Scheme:   "http",
			Host:     "example.com",
			Port:     0,
			Path:     "",
			Query:    "",
			Fragment: "",
		}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			got, err := Of(test.in)

			if !errors.Is(err, test.expectedError) {
				tt.Errorf("Expected error %v, got %v", test.expectedError, err)
			}

			if got != test.want {
				tt.Errorf("Incorrect Uri for %s:\n%s\n", test.in, cmp.Diff(got, test.want))
			}
		})
	}
}
