package shorten_test

import (
	"testing"

	"github.com/Faner201/sc_links/internal/shorten"
	"github.com/stretchr/testify/assert"
)

func TestShorten(t *testing.T) {
	t.Run("the correct short entry is returned", func(t *testing.T) {
		type testCase struct {
			id       uint32
			expected string
		}

		testCases := []testCase{
			{
				id:       40,
				expected: "fk",
			},
			{
				id:       0,
				expected: "",
			},
			{
				id:       80,
				expected: "sk",
			},
		}

		for _, tc := range testCases {
			actual := shorten.Shorten(tc.id)
			assert.Equal(t, tc.expected, actual)
		}
	})
}