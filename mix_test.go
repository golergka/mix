package mix

import "testing"
import "fmt"
import "github.com/stretchr/testify/assert"

func TestAdr(t *testing.T) {
	assert.Equal(t, Adr(0), MakeAdr([]byte{0}))
	assert.Equal(t, Adr(0), MakeAdr([]byte{0, 0}))
	assert.Equal(t, Adr(0), MakeAdr([]byte{0, 0, 0}))
	assert.Equal(t, Adr(0), MakeAdr([]byte{0, 0, 0, 0}))
	assert.Equal(t, Adr(0), MakeAdr([]byte{0, 0, 0, 0, 0}))

	assert.Equal(t, Adr(5), MakeAdr([]byte{5}))
	assert.Equal(t, Adr(64), MakeAdr([]byte{1, 0}))
	assert.Equal(t, Adr(65), MakeAdr([]byte{1, 1}))
	assert.Equal(t, Adr(128), MakeAdr([]byte{2, 0}))
	assert.Equal(t, Adr(325), MakeAdr([]byte{5, 5}))
	assert.Equal(t, Adr(4096), MakeAdr([]byte{1, 0, 0}))
	assert.Equal(t, Adr(4097), MakeAdr([]byte{1, 0, 1}))
}

func TestSignAdr(t *testing.T) {
	var s Sign
	s = false
	assert.Equal(t, SignAdr(-1), MakeSignAdr(&s, []byte{1}))
	assert.Equal(t, SignAdr(-5), MakeSignAdr(&s, []byte{0, 5}))
}

func TestDo(t *testing.T) {
	// LDA
	{
		var m Mix
		m.RI[1] = Index{Sign: true, Bytes: [2]byte{0, 63}}
		m.Memory[31] = Word{Sign: false, Bytes: [5]byte{10, 11, 0, 11, 22}}

		m.Do(&Word{Sign: false, Bytes: [5]byte{0, 32, 02, 11, 8}})

		assert.Equal(t,
			Word{Sign: true, Bytes: [5]byte{00, 00, 10, 11, 00}},
			m.RA)
	}
	// LD3
	{
		var m Mix
		m.RI[0] = Index{Sign: false, Bytes: [2]byte{0, 1}}
		m.RI[2] = Index{Sign: true, Bytes: [2]byte{24, 12}}
		m.Memory[12] = Word{Sign: false, Bytes: [5]byte{1, 2, 3, 4, 5}}

		m.Do(&Word{Sign: true, Bytes: [5]byte{0, 13, 1, 27, 11}})

		assert.Equal(t,
			Index{Sign: true, Bytes: [2]byte{0, 3}},
			m.RI[2])
	}
	// STA
	{
		var m Mix
		originalRA := Word{Sign: true, Bytes: [5]byte{1, 2, 3, 4, 5}}
		m.RA = originalRA
		m.Memory[1200] = Word{Sign: false, Bytes: [5]byte{20, 21, 22, 23, 24}}
		i := "STA 1200(2:3)"
		var w Word
		fmt.Sscanln(i, &w)

		// STA 1200(2:3)
		m.Do(&w)

		assert.Equal(t,
			Word{Sign: false, Bytes: [5]byte{20, 4, 5, 23, 24}},
			m.Memory[1200])

		assert.Equal(t, originalRA, m.RA)
	}
}
