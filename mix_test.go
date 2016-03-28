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

func TestWordField(t * testing.T) {
	assert.Equal( t, 
		  Word{Sign:true,  Bytes:[5]byte{ 0,  0, 10, 11,  0}},
		(&Word{Sign:false, Bytes:[5]byte{10, 11,  0, 11, 22}}).Field(11))
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
	assert.Equal(t, int16(-5), SignedAdr(&s, []byte{0, 5}));
}

func TestWordSignAdr(t *testing.T) {
	assert.Equal(t, int16( 0), (&Word{Sign:false, Bytes:[5]byte{0, 0, 1, 2, 3}}).SignedAdr())
	assert.Equal(t, int16(-1), (&Word{Sign:false, Bytes:[5]byte{0, 1, 1, 2, 3}}).SignedAdr())
	assert.Equal(t, int16(64), (&Word{Sign:true,  Bytes:[5]byte{1, 0, 1, 2, 3}}).SignedAdr())
	assert.Equal(t, int16(-5), (&Word{Sign:false, Bytes:[5]byte{0, 5, 1, 2, 3}}).SignedAdr())
}

func TestWordOpcode(t *testing.T) {
	assert.Equal(t, OP_LDA, (&Word{Bytes:[5]byte{0, 0, 0, 0,  8}}).Opcode())
	assert.Equal(t, OP_LDA, (&Word{Bytes:[5]byte{1, 2, 4, 16, 8}}).Opcode())
}

func TestDo(t *testing.T) {
	{
		var m Mix
		m.RI[1] = Index{Sign:true, Bytes:[2]byte{0, 63}}
		m.Memory[31] = Word {Sign:false, Bytes:[5]byte{10, 11, 0, 11, 22}}
		m.Do(&Word {Sign:false, Bytes:[5]byte{0, 32, 02, 11, 8}})

		assert.Equal(t,
			Word {Sign:true, Bytes:[5]byte{00, 00, 10, 11, 00}},
			m.RA)
	}
	{
		var m Mix
		m.RI[0] = Index{Sign:false, Bytes:[2]byte{0, 1}}
		m.RI[2] = Index{Sign:true, Bytes:[2]byte{24, 12}}
		m.Memory[12] = Word{Sign:false, Bytes:[5]byte{1, 2, 3, 4, 5}}
		m.Do(&Word{Sign:true, Bytes:[5]byte{0, 13, 1, 27, 11}})

		assert.Equal(t,
			Index{Sign:true, Bytes:[2]byte{0, 3}},
			m.RI[2])
	}
}
