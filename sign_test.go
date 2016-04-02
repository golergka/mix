package mix

import "testing"
import "github.com/stretchr/testify/assert"

func TestSignString(t *testing.T) {
	var m Sign
	m = false
	assert.Equal(t, m.String(), "-")
	m = true
	assert.Equal(t, m.String(), "+")
}
