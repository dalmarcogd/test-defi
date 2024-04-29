package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCNPJ struct {
	Document string `validate:"required,cnpj"`
}

func TestCNPJ(t *testing.T) {
	setup, err := Setup()
	assert.NoError(t, err)

	t.Run("invalid format zero", func(t *testing.T) {
		err = setup.Struct(&testCNPJ{Document: "00.000.000/0000-00"})
		assert.EqualError(t, err, "Key: 'testCNPJ.Document' Error:Field validation for 'Document' failed on the 'cnpj' tag")
	})

	t.Run("invalid zero", func(t *testing.T) {
		err = setup.Struct(&testCNPJ{Document: "00000000000000"})
		assert.EqualError(t, err, "Key: 'testCNPJ.Document' Error:Field validation for 'Document' failed on the 'cnpj' tag")
	})

	t.Run("invalid digit format cnpj", func(t *testing.T) {
		err = setup.Struct(&testCNPJ{Document: "25.630.933/0001-94"})
		assert.EqualError(t, err, "Key: 'testCNPJ.Document' Error:Field validation for 'Document' failed on the 'cnpj' tag")
	})

	t.Run("invalid digit cnpj", func(t *testing.T) {
		err = setup.Struct(&testCNPJ{Document: "25630933000194"})
		assert.EqualError(t, err, "Key: 'testCNPJ.Document' Error:Field validation for 'Document' failed on the 'cnpj' tag")
	})

	t.Run("valid format cnpj", func(t *testing.T) {
		err = setup.Struct(&testCNPJ{Document: "25.630.933/0001-93"})
		assert.NoError(t, err)
	})

	t.Run("valid format cnpj", func(t *testing.T) {
		err = setup.Struct(&testCNPJ{Document: "25630933000193"})
		assert.NoError(t, err)
	})
}
