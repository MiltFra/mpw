package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPower(t *testing.T) {
	p := 1
	for i := 0; i < 10; i++ {
		assert.Equal(t, p, GetPower(i))
		p *= 95
	}
}

func TestExtendState(t *testing.T) {
	s := 0
	for i := 0; i < 4; i++ {
		s = ExtendState(4, s, 94)
	}
	s1 := s
	s = ExtendState(4, s, 4)
	for i := 0; i < 4; i++ {
		s = ExtendState(4, s, 94)
	}
	assert.Equal(t, s1, s)
}
