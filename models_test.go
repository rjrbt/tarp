package main

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTarpDetailsMethods(t *testing.T) {
	arbitraryInstance := tarpDetails{
		tarpFunc{
			Filename: "a",
			Name:     "One",
		},
		tarpFunc{
			Filename: "b",
			Name:     "Two",
			DeclPos: token.Position{
				Line: 1,
			},
		},
		tarpFunc{
			Filename: "b",
			Name:     "Three",
			DeclPos: token.Position{
				Line: 2,
			},
		},
	}

	testLen := func(t *testing.T) {
		assert.Equal(t, 3, arbitraryInstance.Len(), ".Len() should return the length of tarpDetails")
	}
	t.Run(".Len()", testLen)

	testLess := func(t *testing.T) {
		assert.True(t, arbitraryInstance.Less(0, 1), ".Less(i, j) should return the correct response")
		assert.False(t, arbitraryInstance.Less(1, 0), ".Less(i, j) should return the correct response")
		assert.True(t, arbitraryInstance.Less(1, 2), ".Less(i, j) should return the correct response")

	}
	t.Run(".Less()", testLess)

	testSwap := func(t *testing.T) {
		expected := tarpDetails{
			tarpFunc{
				Filename: "b",
				Name:     "Two",
				DeclPos: token.Position{
					Line: 1,
				},
			},
			tarpFunc{
				Filename: "a",
				Name:     "One",
			},
			tarpFunc{
				Filename: "b",
				Name:     "Three",
				DeclPos: token.Position{
					Line: 2,
				},
			},
		}
		arbitraryInstance.Swap(0, 1)
		assert.Equal(t, expected, arbitraryInstance, ".Swap(i, j) should swap the location of two values")
	}
	t.Run(".Swap()", testSwap)
}
