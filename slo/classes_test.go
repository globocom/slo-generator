package slo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindClass(t *testing.T) {
	definition := ClassesDefinition{
		Classes: []Class{
			{Name: "HIGH"},
			{Name: "LOW"},
		},
	}

	class, err := definition.FindClass("HIGH")
	assert.NoError(t, err)
	assert.Equal(t, &definition.Classes[0], class)

	class, err = definition.FindClass("NOTFOUND")
	assert.EqualError(t, err, "SLO class \"NOTFOUND\" is not found")
	assert.Nil(t, class)

	class, err = definition.FindClass("")
	assert.Nil(t, err)
	assert.Nil(t, class)
}
