//go:build unit

package validators

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzNotSpaceSuccess(f *testing.F) {
	ctx := context.Background()

	validate, err := Setup()
	assert.NoError(f, err)

	testcases := []string{"1212", "asd0asd", "ASD2eas12345", "myname"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		s := struct {
			Message string `validate:"notspace"`
		}{
			Message: orig,
		}

		assert.NoErrorf(t, validate.StructCtx(ctx, s), "fail to valid value='%s'", orig)
	})
}

func FuzzNotSpaceFail(f *testing.F) {
	ctx := context.Background()

	validate, err := Setup()
	assert.NoError(f, err)

	testcases := []string{" 1212", " asd0asd ", "ASD2e as12345", "my  -name"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		s := struct {
			Message string `validate:"notspace"`
		}{
			Message: orig,
		}

		assert.EqualError(t, validate.StructCtx(ctx, s), "Key: 'Message' Error:Field validation for 'Message' failed on the 'notspace' tag")
	})
}

func TestNotSpaceInvalidType(t *testing.T) {
	ctx := context.Background()

	validate, err := Setup()
	assert.NoError(t, err)

	testcases := []any{"1212 ", 10, int8(1), float64(20), rune(1), byte(123)}
	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("case=%s", testcase), func(t *testing.T) {
			assert.NoError(t, validate.VarCtx(ctx, testcase, "notblank"))
		})
	}
}
