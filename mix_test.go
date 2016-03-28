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

func TestWordString(t *testing.T) {
	var m Word
	assert.Equal(t, "- 00 00 00 00 00", m.String())
}

func TestAdr(t *testing.T) {
	assert.Equal(t, int16(0), Adr([]byte{0}))
	assert.Equal(t, int16(0), Adr([]byte{0, 0}))
	assert.Equal(t, int16(0), Adr([]byte{0, 0, 0}))
	assert.Equal(t, int16(0), Adr([]byte{0, 0, 0, 0}))
	assert.Equal(t, int16(0), Adr([]byte{0, 0, 0, 0, 0}))

	assert.Equal(t, int16(5),		Adr([]byte{5}))
	assert.Equal(t, int16(64),		Adr([]byte{1, 0}))
	assert.Equal(t, int16(65),		Adr([]byte{1, 1}))
	assert.Equal(t, int16(128),		Adr([]byte{2, 0}))
	assert.Equal(t, int16(325),		Adr([]byte{5, 5}))
	assert.Equal(t, int16(4096),	Adr([]byte{1, 0, 0}))
	assert.Equal(t, int16(4097),	Adr([]byte{1, 0, 1}))
}

func TestSignAdr(t *testing.T) {
	var s Sign
	s = false
	assert.Equal(t, int16(-1), SignedAdr(&s, []byte{1}));
}
