package format_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rockset/cli/format"
)

func TestParseSelectionString(t *testing.T) {
	tests := []string{
		"Foo:.foo:size",
		"Foo:.foo,Baz:.bar.baz:size",
	}

	for _, tst := range tests {
		s, err := format.ParseSelectionString(tst)

		require.NoError(t, err)
		assert.Equal(t, tst, s.String())
	}
}
