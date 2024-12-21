package markup_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ci-space/edit-config/internal/shared/markup"
)

func TestYAMLDocument_Append(t *testing.T) {
	cases := []struct {
		Title string

		Input    string
		Pointer  string
		Value    interface{}
		Expected string
	}{
		{
			Title: "append string to string",

			Input:    "foo: bar",
			Pointer:  "foo",
			Value:    "new",
			Expected: "foo: barnew\n",
		},
		{
			Title: "append int to int",

			Input:    "foo: 1",
			Pointer:  "foo",
			Value:    2,
			Expected: "foo: 3\n",
		},
		{
			Title: "append string to slice of strings",

			Input:    "foo: [a, b, c]",
			Pointer:  "foo",
			Value:    "d",
			Expected: "foo: [a, b, c, d]\n",
		},
		{
			Title: "append int to slice of ints",

			Input:    "foo: [1, 2, 3]",
			Pointer:  "foo",
			Value:    4,
			Expected: "foo: [1, 2, 3, 4]\n",
		},
		{
			Title: "append numeric string to slice of ints",

			Input:    "foo: [1, 2, 3]",
			Pointer:  "foo",
			Value:    "4",
			Expected: "foo: [1, 2, 3, 4]\n",
		},
		{
			Title: "append bool to slice of bools",

			Input:    "foo: [true, false, false]",
			Pointer:  "foo",
			Value:    true,
			Expected: "foo: [true, false, false, true]\n",
		},
		{
			Title: "append bool string to slice of bools",

			Input:    "foo: [true, false, false]",
			Pointer:  "foo",
			Value:    "true",
			Expected: "foo: [true, false, false, true]\n",
		},
		{
			Title: "append float64 to slice of float64s",

			Input:    "foo: [10.5, 10.6, 10.7]",
			Pointer:  "foo",
			Value:    10.8,
			Expected: "foo: [10.5, 10.6, 10.7, 10.8]\n",
		},
		{
			Title: "append float64 string to slice of float64s",

			Input:    "foo: [10.5, 10.6, 10.7]",
			Pointer:  "foo",
			Value:    "10.8",
			Expected: "foo: [10.5, 10.6, 10.7, 10.8]\n",
		},
	}

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			doc, err := markup.LoadYAMLDocument([]byte(c.Input))
			require.NoError(t, err)

			err = doc.Append(c.Pointer, c.Value)
			require.NoError(t, err)

			assert.Equal(t, c.Expected, doc.String())
		})
	}
}
