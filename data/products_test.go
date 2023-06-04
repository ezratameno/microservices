package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksValidation(t *testing.T) {
	p := &Product{}

	assert.Error(t, p.Validate())

	p = &Product{
		Name:  "sdfsd",
		Price: 415,
		SKU:   "sdf-sdf-sd",
	}
	err := p.Validate()

	fmt.Println(err)
	assert.Nil(t, err)

}
