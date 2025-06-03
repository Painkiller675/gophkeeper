package version

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPrintBuildInfo is a simple test for version setting
func TestPrintBuildInfo(t *testing.T) {
	for _, tc := range []struct {
		name     string
		version  string
		date     string
		author   string
		commit   string
		expected string
	}{
		{
			name:    "Simple test",
			version: "0.0.1", date: "2025/06/02", author: "Painkiller675", commit: "12345678",
			expected: "0.0.1, date: 2025/06/02, author: Painkiller675, commit: 12345678\n",
		},
		{
			name:     "Empty fields",
			expected: "N/A, date: N/A, commit: N/A\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var builder strings.Builder

			BuildVersion = tc.version
			BuildDate = tc.date
			BuildCommit = tc.commit

			WriteBuildInfo(&builder)
			assert.Equal(t, tc.expected, builder.String())
			assert.Equal(t, tc.expected, Info())
		})
	}
}
