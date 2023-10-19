package format_test

import (
	"github.com/rockset/rockset-go-client/openapi"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rockset/cli/format"
)

func TestParseSelectionString(t *testing.T) {
	tests := []string{
		"Foo:.foo:size",
		"Foo:.foo,Baz:.bar.baz:size",
		"Arr:.foo.arr[1]",
		"ArrMap:.foo.arr[].value",
	}

	for _, tst := range tests {
		s, err := format.ParseSelectionString(tst)

		require.NoError(t, err)
		assert.Equal(t, tst, s.String())
	}
}

func TestBasicSelection(t *testing.T) {
	hello := "hello"
	collection := openapi.Collection{Name: &hello}

	s, err := format.ParseSelectionString(".name")
	require.NoError(t, err)

	name, err := s[0].Select(collection)
	require.NoError(t, err)

	assert.Equal(t, hello, name)
}

func TestSelectArrayMapping(t *testing.T) {
	hello := "hello"
	world := "world"

	integration := openapi.Integration{
		Collections: []openapi.Collection{{Name: &hello}, {Name: &world}},
	}

	s, err := format.ParseSelectionString(".collections[].name")
	require.NoError(t, err)

	arr, err := s[0].Select(integration)

	switch typedArr := arr.(type) {
	case []any:
		assert.Equal(t, 2, len(typedArr))
		assert.Equal(t, "hello", typedArr[0])
		assert.Equal(t, "world", typedArr[1])
		break
	default:
		assert.Fail(t, "returned non-array type")
	}

}
