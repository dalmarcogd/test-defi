package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCPF struct {
	Document string `validate:"required,cpf"`
}

func TestCPF(t *testing.T) {
	setup, err := Setup()
	assert.NoError(t, err)

	t.Run("invalid format zero", func(t *testing.T) {
		err = setup.Struct(&testCPF{Document: "000.000.000-00"})
		assert.EqualError(t, err, "Key: 'testCPF.Document' Error:Field validation for 'Document' failed on the 'cpf' tag")
	})

	t.Run("invalid zero", func(t *testing.T) {
		err = setup.Struct(&testCPF{Document: "00000000000"})
		assert.EqualError(t, err, "Key: 'testCPF.Document' Error:Field validation for 'Document' failed on the 'cpf' tag")
	})

	t.Run("invalid digit format cpf", func(t *testing.T) {
		err = setup.Struct(&testCPF{Document: "831.110.680-04"})
		assert.EqualError(t, err, "Key: 'testCPF.Document' Error:Field validation for 'Document' failed on the 'cpf' tag")
	})

	t.Run("invalid digit cpf", func(t *testing.T) {
		err = setup.Struct(&testCPF{Document: "83111068004"})
		assert.EqualError(t, err, "Key: 'testCPF.Document' Error:Field validation for 'Document' failed on the 'cpf' tag")
	})

	t.Run("valid format cpf", func(t *testing.T) {
		err = setup.Struct(&testCPF{Document: "831.110.680-05"})
		assert.NoError(t, err)
	})

	t.Run("valid format cpf", func(t *testing.T) {
		err = setup.Struct(&testCPF{Document: "83111068005"})
		assert.NoError(t, err)
	})
}
