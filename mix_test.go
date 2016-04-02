package mix

import "testing"
import "fmt"
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

func TestWordFormat(t *testing.T) {
	{
		m := Word{}

		s := fmt.Sprintf("%v", m)

		assert.Equal(t, "- 00 00 00 00 00", s)
	}
	{
		m := Word{Sign:true, Bytes:[5]byte{18, 48, 0, 19, 24}}

		s := fmt.Sprintf("%i", m)

		assert.Equal(t, "STA 1200(2:3)", s)
	}
	{
		m := Word{Sign:false, Bytes:[5]byte{0, 32, 2, 11, 10}}

		s := fmt.Sprintf("%i", m)

		assert.Equal(t, "LD2 -32,2(1:3)", s)
	}
}

func TestWordScan(t *testing.T) {
	{
		var m Word
		var w interface{}

		w = &m

		if _, ok := w.(fmt.Scanner); !ok {
			t.Error("Word doesn't imlpement Scanner")
		}
	}
	{
		s := "- 00 00 00 00 00"
		var m Word

		_, err := fmt.Sscanf(s, "%v", &m)

		assert.Nil(t, err)
		assert.Equal(t, Word{}, m)
	}
	{
		s := "+ 1 2 3 4 5"
		var m Word

		_, err := fmt.Sscanf(s, "%v", &m)

		assert.Nil(t, err)
		assert.Equal(t, Word{Sign:true, Bytes:[5]byte{1, 2, 3, 4, 5}}, m)
	}
}

func TestOpScan(t *testing.T) {
	{
		var o Op
		s := "LDA"

		_, err := fmt.Sscanf(s, "%v", &o)

		assert.Nil(t, err)
		assert.Equal(t, LDA, o)
	}
	{
		var o Op
		s := "asdasd"
		
		_, err := fmt.Sscanf(s, "%v", &o)

		assert.NotNil(t, err)
	}
}

func TestWordField(t * testing.T) {
	assert.Equal( t, 
		  Word{Sign:true,  Bytes:[5]byte{ 0,  0, 10, 11,  0}},
		(&Word{Sign:false, Bytes:[5]byte{10, 11,  0, 11, 22}}).GetField(11))
}

func TestAdr(t *testing.T) {
	assert.Equal(t, Adr(0), MakeAdr([]byte{0}))
	assert.Equal(t, Adr(0), MakeAdr([]byte{0, 0}))
	assert.Equal(t, Adr(0), MakeAdr([]byte{0, 0, 0}))
	assert.Equal(t, Adr(0), MakeAdr([]byte{0, 0, 0, 0}))
	assert.Equal(t, Adr(0), MakeAdr([]byte{0, 0, 0, 0, 0}))

	assert.Equal(t, Adr(5),		MakeAdr([]byte{5}))
	assert.Equal(t, Adr(64),	MakeAdr([]byte{1, 0}))
	assert.Equal(t, Adr(65),	MakeAdr([]byte{1, 1}))
	assert.Equal(t, Adr(128),	MakeAdr([]byte{2, 0}))
	assert.Equal(t, Adr(325),	MakeAdr([]byte{5, 5}))
	assert.Equal(t, Adr(4096),	MakeAdr([]byte{1, 0, 0}))
	assert.Equal(t, Adr(4097),	MakeAdr([]byte{1, 0, 1}))
}

func TestSignAdr(t *testing.T) {
	var s Sign
	s = false
	assert.Equal(t, SignAdr(-1), MakeSignAdr(&s, []byte{1}));
	assert.Equal(t, SignAdr(-5), MakeSignAdr(&s, []byte{0, 5}));
}

func TestWordSignAdr(t *testing.T) {
	assert.Equal(t, SignAdr( 0), (&Word{Sign:false, Bytes:[5]byte{0, 0, 1, 2, 3}}).GetAdr())
	assert.Equal(t, SignAdr(-1), (&Word{Sign:false, Bytes:[5]byte{0, 1, 1, 2, 3}}).GetAdr())
	assert.Equal(t, SignAdr(64), (&Word{Sign:true,  Bytes:[5]byte{1, 0, 1, 2, 3}}).GetAdr())
	assert.Equal(t, SignAdr(-5), (&Word{Sign:false, Bytes:[5]byte{0, 5, 1, 2, 3}}).GetAdr())
}

func TestWordGetOp(t *testing.T) {
	assert.Equal(t, LDA, (&Word{Bytes:[5]byte{0, 0, 0, 0,  8}}).GetOp())
	assert.Equal(t, LDA, (&Word{Bytes:[5]byte{1, 2, 4, 16, 8}}).GetOp())
}

func TestWordSetOp(t *testing.T) {
	{
		w := Word{Bytes:[5]byte{1, 2, 3, 4, 5}}

		w.SetOp(LDA)

		assert.Equal(t, byte(8), w.Bytes[4])
	}
}

func TestDo(t *testing.T) {
	// LDA
	{
		var m Mix
		m.RI[1] = Index{Sign:true, Bytes:[2]byte{0, 63}}
		m.Memory[31] = Word {Sign:false, Bytes:[5]byte{10, 11, 0, 11, 22}}

		m.Do(&Word {Sign:false, Bytes:[5]byte{0, 32, 02, 11, 8}})

		assert.Equal(t,
			Word {Sign:true, Bytes:[5]byte{00, 00, 10, 11, 00}},
			m.RA)
	}
	// LD3
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
	// STA
	{
		var m Mix
		originalRA := Word{Sign:true, Bytes:[5]byte{1, 2, 3, 4, 5}}
		m.RA = originalRA
		m.Memory[1200] = Word{Sign:false, Bytes:[5]byte{20, 21, 22, 23, 24}}
		
		// STA 1200(2:3)
		m.Do(&Word{Sign:true, Bytes:[5]byte{18, 48, 0, 19, 24}})

		assert.Equal(t,
			Word{Sign:false, Bytes:[5]byte{20, 4, 5, 23, 24}},
			m.Memory[1200])

		assert.Equal(t, originalRA, m.RA)
	}
}
