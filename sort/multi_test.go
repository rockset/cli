package sort_test

import (
	"github.com/rockset/cli/sort"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStruct struct {
	X string
	Y int
}

func TestSort(t *testing.T) {
	list := []TestStruct{
		{X: "a", Y: 2},
		{X: "b", Y: 2},
		{X: "b", Y: 1},
		{X: "a", Y: 1},
	}

	m := sort.Multi[TestStruct]{LessFuncs: []func(*TestStruct, *TestStruct) bool{
		func(testStruct *TestStruct, testStruct2 *TestStruct) bool {
			return testStruct.X < testStruct2.X
		},
		func(testStruct *TestStruct, testStruct2 *TestStruct) bool {
			return testStruct.Y < testStruct2.Y
		},
	}}
	m.Sort(list)

	assert.Equal(t, "a", list[0].X)
	assert.Equal(t, 1, list[0].Y)
	assert.Equal(t, "a", list[1].X)
	assert.Equal(t, 2, list[1].Y)
	assert.Equal(t, "b", list[2].X)
	assert.Equal(t, 1, list[2].Y)
	assert.Equal(t, "b", list[3].X)
	assert.Equal(t, 2, list[3].Y)
}
